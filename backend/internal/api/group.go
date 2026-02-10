package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"badminton_tournament/backend/internal/models"
)

type CreateGroupRequest struct {
	Name          string      `json:"name"`
	Pool          string      `json:"pool"` // "Mesoneer" or "Lab"
	TournamentID  uuid.UUID   `json:"tournament_id"`
	TeamIDs       []uuid.UUID `json:"team_ids"` // Expect exactly 4 IDs
}

func (h *Handler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.TeamIDs) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group must have exactly 4 teams"})
		return
	}
	if req.Pool == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool is required (Mesoneer or Lab)"})
		return
	}

	ctx := c.Request.Context()

	// 1. Validate Teams: Must be in same Pool and not busy
	// Fetch teams to check pool
	var teams []models.Team
	if err := h.DB.NewSelect().Model(&teams).Where("id IN (?)", bun.In(req.TeamIDs)).Scan(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
		return
	}
	
	if len(teams) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more teams not found"})
		return
	}

	for _, team := range teams {
		if team.Pool != req.Pool {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All teams must belong to the selected Pool (" + req.Pool + ")"})
			return
		}
	}

	// Shuffle Teams for Random Seeding
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(req.TeamIDs), func(i, j int) {
		req.TeamIDs[i], req.TeamIDs[j] = req.TeamIDs[j], req.TeamIDs[i]
	})

	// Check if already in active match
	count, err := h.DB.NewSelect().Model((*models.Match)(nil)).
		Where("team_a_id IN (?) OR team_b_id IN (?)", bun.In(req.TeamIDs), bun.In(req.TeamIDs)).
		Count(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate team availability"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more selected teams are already competing in another group"})
		return
	}

	// 2. Create Group
	group := &models.Group{
		TournamentID: req.TournamentID,
		Name:         req.Name,
		Pool:         req.Pool,
	}
	_, err = h.DB.NewInsert().Model(group).Returning("*").Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// 2. Create 5 Matches (GSL Structure)
	// M5 (Decider)
	m5 := &models.Match{GroupID: group.ID, Label: "Decider"} // Match 5
	_, _ = h.DB.NewInsert().Model(m5).Exec(ctx)

	// M3 (Winners Match) -> Win: Qualify, Lose: Go to M5
	m3 := &models.Match{GroupID: group.ID, Label: "Winners", NextMatchLoseID: m5.ID}
	_, _ = h.DB.NewInsert().Model(m3).Exec(ctx)
	
	// M4 (Losers Match) -> Win: Go to M5, Lose: Out
	m4 := &models.Match{GroupID: group.ID, Label: "Losers", NextMatchWinID: m5.ID}
	_, _ = h.DB.NewInsert().Model(m4).Exec(ctx)

	// M1 (Opening A) -> Win: M3, Lose: M4
	m1 := &models.Match{
		GroupID:         group.ID,
		Label:           "M1",
		TeamAID:         req.TeamIDs[0],
		TeamBID:         req.TeamIDs[1],
		NextMatchWinID:  m3.ID,
		NextMatchLoseID: m4.ID,
	}
	_, _ = h.DB.NewInsert().Model(m1).Exec(ctx)

	// M2 (Opening B) -> Win: M3, Lose: M4
	m2 := &models.Match{
		GroupID:         group.ID,
		Label:           "M2",
		TeamAID:         req.TeamIDs[2],
		TeamBID:         req.TeamIDs[3],
		NextMatchWinID:  m3.ID,
		NextMatchLoseID: m4.ID,
	}
	_, err = h.DB.NewInsert().Model(m2).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create matches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group_id": group.ID, "status": "created"})
}

func (h *Handler) ListGroups(c *gin.Context) {
	var groups []models.Group
	
	// Eager load matches and teams
	err := h.DB.NewSelect().Model(&groups).
		Relation("Matches", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("label ASC").
				Relation("TeamA").
				Relation("TeamB").
				Relation("Winner")
		}).
		Order("name ASC").
		Scan(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"badminton_tournament/backend/internal/models"
)

type CreateGroupRequest struct {
	Name         string      `json:"name"`
	Pool         string      `json:"pool"` // "Mesoneer" or "Lab"
	TournamentID uuid.UUID   `json:"tournament_id"`
	TeamIDs      []uuid.UUID `json:"team_ids"` // Expect exactly 4 IDs
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	// 3. Create GSL Matches
	if err := h.createGSLMatches(ctx, group.ID, req.TeamIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create matches: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group_id": group.ID, "status": "created"})
}

type AutoGenerateGroupsRequest struct {
	Pool         string    `json:"pool"`
	TournamentID uuid.UUID `json:"tournament_id"`
	NamePrefix   string    `json:"name_prefix"`
}

func (h *Handler) AutoGenerateGroups(c *gin.Context) {
	var req AutoGenerateGroupsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Pool == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool is required"})
		return
	}

	ctx := c.Request.Context()

	// Fetch available teams in pool
	var availableTeams []models.Team
	err := h.DB.NewSelect().
		Model(&availableTeams).
		Where("pool = ?", req.Pool).
		Where("id NOT IN (SELECT team_a_id FROM matches WHERE team_a_id IS NOT NULL UNION SELECT team_b_id FROM matches WHERE team_b_id IS NOT NULL)").
		Scan(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available teams: " + err.Error()})
		return
	}

	numTeams := len(availableTeams)
	if numTeams == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available teams in " + req.Pool + " pool"})
		return
	}

	if numTeams%4 != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot auto-generate: " + fmt.Sprintf("%d", numTeams) + " teams available, but groups must have exactly 4 teams"})
		return
	}

	// Shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numTeams, func(i, j int) {
		availableTeams[i], availableTeams[j] = availableTeams[j], availableTeams[i]
	})

	var createdGroups []uuid.UUID
	for i := 0; i < numTeams; i += 4 {
		name := req.NamePrefix
		if name == "" {
			name = "Group"
		}
		name = fmt.Sprintf("%s %d", name, (i/4)+1)

		group := &models.Group{
			TournamentID: req.TournamentID,
			Name:         name,
			Pool:         req.Pool,
		}
		_, err := h.DB.NewInsert().Model(group).Returning("*").Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
			return
		}

		teamIDs := []uuid.UUID{
			availableTeams[i].ID,
			availableTeams[i+1].ID,
			availableTeams[i+2].ID,
			availableTeams[i+3].ID,
		}

		if err := h.createGSLMatches(ctx, group.ID, teamIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create matches for group " + name + ": " + err.Error()})
			return
		}
		createdGroups = append(createdGroups, group.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "created",
		"groups_created": len(createdGroups),
		"group_ids":      createdGroups,
	})
}

func (h *Handler) createGSLMatches(ctx context.Context, groupID uuid.UUID, teamIDs []uuid.UUID) error {
	// M5 (Decider)
	m5 := &models.Match{GroupID: groupID, Label: "Decider"}
	_, err := h.DB.NewInsert().Model(m5).Exec(ctx)
	if err != nil {
		return err
	}

	// M3 (Winners Match)
	m3 := &models.Match{GroupID: groupID, Label: "Winners", NextMatchLoseID: m5.ID}
	_, err = h.DB.NewInsert().Model(m3).Exec(ctx)
	if err != nil {
		return err
	}

	// M4 (Losers Match)
	m4 := &models.Match{GroupID: groupID, Label: "Losers", NextMatchWinID: m5.ID}
	_, err = h.DB.NewInsert().Model(m4).Exec(ctx)
	if err != nil {
		return err
	}

	// M1 (Opening A)
	m1 := &models.Match{
		GroupID:         groupID,
		Label:           "M1",
		TeamAID:         teamIDs[0],
		TeamBID:         teamIDs[1],
		NextMatchWinID:  m3.ID,
		NextMatchLoseID: m4.ID,
	}
	_, err = h.DB.NewInsert().Model(m1).Exec(ctx)
	if err != nil {
		return err
	}

	// M2 (Opening B)
	m2 := &models.Match{
		GroupID:         groupID,
		Label:           "M2",
		TeamAID:         teamIDs[2],
		TeamBID:         teamIDs[3],
		NextMatchWinID:  m3.ID,
		NextMatchLoseID: m4.ID,
	}
	_, err = h.DB.NewInsert().Model(m2).Exec(ctx)
	return err
}

func (h *Handler) ListGroups(c *gin.Context) {
	var groups []models.Group

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

package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Pool string

const (
	PoolMesoneer Pool = "Mesoneer"
	PoolLab      Pool = "Lab"
)

type Participant struct {
	bun.BaseModel `bun:"table:participants,alias:p"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Pool      Pool      `bun:"pool,notnull" json:"pool"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"` // E.g. "Player1 & Player2"
	Player1ID int64     `bun:"player1_id,notnull" json:"player1_id"`
	Player2ID int64     `bun:"player2_id,notnull" json:"player2_id"`
	Pool      Pool      `bun:"pool,notnull" json:"pool"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Match struct {
	bun.BaseModel `bun:"table:matches,alias:m"`

	ID          int64  `bun:"id,pk,autoincrement" json:"id"`
	Description string `bun:"description" json:"description"` // e.g. "Group A - Match 1"
	Pool        Pool   `bun:"pool,notnull" json:"pool"`
	DataGroup   string `bun:"group_name" json:"group_name"` // "Group A"

	TeamAID int64 `bun:"team_a_id,nullzero" json:"team_a_id"` // Can be null if TBD
	TeamA   *Team `bun:"rel:belongs-to,join:team_a_id=id" json:"team_a"`
	TeamBID int64 `bun:"team_b_id,nullzero" json:"team_b_id"` // Can be null if TBD
	TeamB   *Team `bun:"rel:belongs-to,join:team_b_id=id" json:"team_b"`

	ScoreA   int `bun:"score_a,default:0" json:"score_a"`
	ScoreB   int `bun:"score_b,default:0" json:"score_b"`
	WinnerID int64 `bun:"winner_id,nullzero" json:"winner_id"` // Set when finished

	VideoURL string `bun:"video_url" json:"video_url"`

	// Linkage for automation
	NextMatchID     int64  `bun:"next_match_id,nullzero" json:"next_match_id"`
	NextMatchSlot   string `bun:"next_match_slot" json:"next_match_slot"` // "A" or "B"
	LoserNextMatchID     int64  `bun:"loser_next_match_id,nullzero" json:"loser_next_match_id"`
	LoserNextMatchSlot   string `bun:"loser_next_match_slot" json:"loser_next_match_slot"` // "A" or "B"

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

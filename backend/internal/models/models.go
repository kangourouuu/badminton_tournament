package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tournament struct {
	bun.BaseModel `bun:"table:tournaments,alias:t"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Status    string    `bun:"status,notnull" json:"status"` // 'draft', 'active', 'completed'
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Participant struct {
	bun.BaseModel `bun:"table:participants,alias:p"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Name      string    `bun:"name,notnull" json:"name"`
	Pool      string    `bun:"pool,notnull" json:"pool"` // 'Mesoneer', 'Lab'
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:tm"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Player1ID uuid.UUID `bun:"player1_id,type:uuid,notnull" json:"player1_id"`
	Player2ID uuid.UUID `bun:"player2_id,type:uuid" json:"player2_id"` // Nullable logic handled by pointer or omitted if strict
	Pool      string    `bun:"pool,notnull" json:"pool"`
	Name      string    `bun:"name,notnull" json:"name"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	// Relations
	Player1 *Participant `bun:"rel:belongs-to,join:player1_id=id" json:"player1,omitempty"`
	Player2 *Participant `bun:"rel:belongs-to,join:player2_id=id" json:"player2,omitempty"`
}

type Group struct {
	bun.BaseModel `bun:"table:groups,alias:g"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	TournamentID uuid.UUID `bun:"tournament_id,type:uuid" json:"tournament_id"`
	Name         string    `bun:"name,notnull" json:"name"` // "Group A"

	// Relations
	Matches []*Match `bun:"rel:has-many,join:id=group_id" json:"matches,omitempty"`
}

type Match struct {
	bun.BaseModel `bun:"table:matches,alias:m"`

	ID      uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	GroupID uuid.UUID `bun:"group_id,type:uuid" json:"group_id"`
	Label   string    `bun:"label,notnull" json:"label"` // "M1", "M2", "Winners", "Losers", "Decider"

	TeamAID uuid.UUID `bun:"team_a_id,type:uuid,nullzero" json:"team_a_id,omitempty"`
	TeamBID uuid.UUID `bun:"team_b_id,type:uuid,nullzero" json:"team_b_id,omitempty"`

	WinnerID uuid.UUID `bun:"winner_id,type:uuid,nullzero" json:"winner_id,omitempty"`
	Score    string    `bun:"score" json:"score"`         // "21-19, 21-18"
	VideoURL string    `bun:"video_url" json:"video_url"` // YouTube link

	// Automation Linking
	NextMatchWinID  uuid.UUID `bun:"next_match_win_id,type:uuid,nullzero" json:"next_match_win_id,omitempty"`
	NextMatchLoseID uuid.UUID `bun:"next_match_lose_id,type:uuid,nullzero" json:"next_match_lose_id,omitempty"`

	// Relations
	TeamA  *Team `bun:"rel:belongs-to,join:team_a_id=id" json:"team_a,omitempty"`
	TeamB  *Team `bun:"rel:belongs-to,join:team_b_id=id" json:"team_b,omitempty"`
	Winner *Team `bun:"rel:belongs-to,join:winner_id=id" json:"winner,omitempty"`
}

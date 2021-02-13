package party

import (
	"time"
)

// Info represents a single meet.
type Info struct {
	ID          string    `db:"party_id"`
	Name        string    `db:"name"`
	Location    string    `db:"location"`
	LfPlayers   int       `db:"lf_players"`
	LfGM        int       `db:"lf_gm"`
	Description string    `db:"description"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

// NewParty describes the required data for creating a new party.
type NewParty struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"required,min=50,max=1000"`
	Location    string `json:"location" validate:"required,max=255"`
	LfPlayers   int    `json:"lf_players,string" validate:"required_without=LfGM,gte=0,lte=10"`
	LfGM        int    `json:"lf_gm,string" validate:"required_without=LfPlayers,gte=0,lte=10"`
}

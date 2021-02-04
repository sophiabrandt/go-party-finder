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
	Title       string
	Description string
	Location    string
	LfPlayers   int
	LfGM        int
}

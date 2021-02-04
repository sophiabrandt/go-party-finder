// Package party contains party-related CRUD functionality
package party

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/internal/database"
)

var (
	// ErrNotFound is used when a specific Party is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Party manages the set of API's for party access.
type Party struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Party for api access.
func New(log *log.Logger, db *sqlx.DB) Party {
	return Party{
		log: log,
		db:  db,
	}
}

// Query gets all parties from the database.
func (p Party) Query(ctx context.Context, traceID string, pageNumber int, rowsPerPage int) ([]*Info, error) {
	const q = `
	SELECT
		party_id, name, location, description, lf_players, lf_gm
	FROM
		parties AS p
	ORDER BY
		party_id
	OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	p.log.Printf("%s: %s: %s", traceID, "party.Query",
		database.Log(q, offset, rowsPerPage),
	)

	parties := []*Info{}
	if err := p.db.SelectContext(ctx, &parties, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting parties")
	}

	return parties, nil
}

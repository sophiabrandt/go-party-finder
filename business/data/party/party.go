// Package party contains party-related CRUD functionality
package party

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/foundation/database"
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

// QueryByID finds the party identified by a given ID.
func (p Party) QueryByID(ctx context.Context, traceID string, partyID string) (*Info, error) {
	// intialize pointer to a zeroed Info struct
	prty := &Info{}

	if _, err := uuid.Parse(partyID); err != nil {
		return prty, ErrInvalidID
	}

	const q = `
	SELECT
		party_id, name, location, description, lf_players, lf_gm
	FROM
		parties AS p
	WHERE
		p.party_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "party.QueryByID",
		database.Log(q, partyID),
	)

	if err := p.db.GetContext(ctx, prty, q, partyID); err != nil {
		if err == sql.ErrNoRows {
			return prty, ErrNotFound
		}
		return prty, errors.Wrap(err, "selecting single product")
	}

	return prty, nil
}

// Create adds a Party to the database. It returns the created Party with
// fields like ID and DateCreated populated.
func (p Party) Create(ctx context.Context, traceID string, np NewParty, now time.Time) (Info, error) {
	prty := Info{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Location:    np.Location,
		Description: np.Description,
		LfPlayers:   np.LfPlayers,
		LfGM:        np.LfGM,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
	INSERT INTO parties
		(party_id, name, location, description, lf_players, lf_gm, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)`

	p.log.Printf("%s: %s: %s", traceID, "party.Create",
		database.Log(q, prty.ID, prty.Name, prty.Location, prty.Description, prty.LfPlayers, prty.LfGM, prty.DateCreated, prty.DateUpdated),
	)

	if _, err := p.db.ExecContext(ctx, q, prty.ID, prty.Name, prty.Location, prty.Description, prty.LfPlayers, prty.LfGM, prty.DateCreated, prty.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting prty")
	}

	return prty, nil
}

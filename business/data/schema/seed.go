package schema

import (
	"github.com/jmoiron/sqlx"
)

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.
const seeds = `
-- Create parties
INSERT INTO parties (party_id, name, location, description, lf_players, lf_gm, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Cologne D&D Group', 'Cologne', 'Ipsa aliquam eaque quo consequuntur adipisci distinctio quam. Nostrum quia atque voluptatibus omnis dolorum. Velit ex autem magni officia accusantium nihil. Earum nulla quam nostrum doloribus quae inventore ipsum accusantium.', 3, 0, '2021-01-01 00:00:01.000001+00', '2020-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'Basic Fantasy Saturdays', 'Frankfurt', 'Iure laboriosam nobis optio magni sint optio. Et unde quos quae minus alias. Beatae sint quo laboriosam.', 1, 1, '2021-01-01 00:00:01.000001+00', '2020-01-01 00:00:01.000001+00'),
	('68deba86-06ed-4397-82d1-98904816acd2', 'Story Gamers Berlin', 'Berlin', 'Eos ullam non tenetur voluptas qui. Eligendi et ut exercitationem soluta. Ut veritatis voluptatem assumenda pariatur qui ex quod. Harum aliquam eum doloremque enim sint omnis labore repudiandae.', 0, 0, '2021-01-01 00:00:01.000001+00', '2020-01-01 00:00:01.000001+00')
	ON CONFLICT DO NOTHING;
`

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM parties;`

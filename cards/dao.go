// Cards DAO.
//
// Rémy Mathieu © 2016

package cards

import (
	"database/sql"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
)

type CardsDAO struct {
	DB *sql.DB
}

// ----------------------

var dao *CardsDAO

func DAO() *CardsDAO {
	if dao != nil {
		return dao
	}

	dao = &CardsDAO{
		DB: storage.DB(),
	}

	if err := dao.InitStmt(); err != nil {
		log.Error("Can't prepare CardsDAO")
		panic(err)
	}

	return dao
}

func (d *CardsDAO) InitStmt() error {
	var err error
	return err
}

// GetByUser returns the cards of the given user.
func (d *CardsDAO) GetByUser(uid string, state CardState) ([]SimpleCard, error) {
	rv := make([]SimpleCard, 0)

	rows, err := d.DB.Query(`
		SELECT "uid", "text", "position"
		FROM "card"
		WHERE
			"user_uid" = $1
			AND
			"state" = $2
	`, uid, state.String())

	if err != nil || rows == nil {
		return rv, err
	}

	defer rows.Close()
	for rows.Next() {
		var sc SimpleCard

		if err := rows.Scan(&sc.Uid, &sc.Text, &sc.Position); err != nil {
			return rv, err
		}

		rv = append(rv, sc)
	}

	return rv, nil
}

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

type CardState string

const (
	// CardActive is an active card of the   user.
	CardActive = "CardActive"
	// CardArchived has been archived by the user.
	CardArchived = "CardArchived"
)

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

	return rv, nil
}

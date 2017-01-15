// Account DAO.
//
// Rémy Mathieu © 2016

package cards

import (
	"database/sql"
	"fmt"
	"time"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"
)

type AccountDAO struct {
	DB *sql.DB
}

// ----------------------

var dao *AccountDAO

func DAO() *AccountDAO {
	if dao != nil {
		return dao
	}

	dao = &AccountDAO{
		DB: storage.DB(),
	}

	if err := dao.InitStmt(); err != nil {
		log.Error("Can't prepare AccountDAO")
		panic(err)
	}

	return dao
}

func (d *AccountDAO) InitStmt() error {
	var err error
	return err
}

// Create inserts the given account in database.
func (d *AccountDAO) Create(uid uuid.UUID, firstname, email, hash string, t time.Time) error {
	var err error
	var r *sql.Result

	if r, err = d.DB.Exec(`
		INSERT INTO "user"
		("uid", "email", "firstname", "hash", "creation_time", "last_update")
		VALUES
		($1, $2, $3, $4, $5, $6)
	`, uid, email, firstname, hash, t, t); err != nil {
		return log.Err("Account/Create:", err)
	}

	return nil
}

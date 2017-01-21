// Account DAO.
//
// Rémy Mathieu © 2016

package accounts

import (
	"database/sql"
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

// UidByEmail returns the uid attached to the given
// email if it is already used.
// Otherwise, returns nil.
func (d *AccountDAO) UidByEmail(email string) (uuid.UUID, error) {
	var err error
	var uid uuid.UUID

	if d.DB.QueryRow(`
		SELECT "uid"
		FROM "user"
		WHERE
			"email" = $1
	`, email).Scan(&uid); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return uid, nil
}

// UserByUid returns basic information of an user
// using its uid.
func (d *AccountDAO) UserByUid(uid uuid.UUID) (SimpleUser, string, error) {
	if uid.IsNil() {
		return SimpleUser{}, "", nil
	}

	var err error
	var su SimpleUser
	var hash string

	if d.DB.QueryRow(`
		SELECT "uid", "firstname", "email", "hash"
		FROM "user"
		WHERE
			"uid" = $1
	`, uid).Scan(&su.Uid, &su.Firstname, &su.Email, &hash); err != nil && err != sql.ErrNoRows {
		return su, "", err
	}

	return su, hash, nil
}

// UserByEmail returns basic information of an user
// using its email as unique identifier.
func (d *AccountDAO) UserByEmail(email string) (SimpleUser, string, error) {
	var err error
	var su SimpleUser
	var hash string

	if d.DB.QueryRow(`
		SELECT "uid", "firstname", "email", "hash"
		FROM "user"
		WHERE
			"email" = $1
	`, email).Scan(&su.Uid, &su.Firstname, &su.Email, &hash); err != nil && err != sql.ErrNoRows {
		return su, "", err
	}

	return su, hash, nil
}

// Create inserts the given account in database.
func (d *AccountDAO) Create(uid uuid.UUID, firstname, email, hash string, t time.Time) error {
	var err error

	if _, err = d.DB.Exec(`
		INSERT INTO "user"
		("uid", "email", "firstname", "hash", "creation_time", "last_update")
		VALUES
		($1, $2, $3, $4, $5, $6)
	`, uid, email, firstname, hash, t, t); err != nil {
		return log.Err("Account/Create:", err)
	}

	return nil
}

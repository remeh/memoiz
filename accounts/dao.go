// Account DAO.
//
// Rémy Mathieu © 2016

package accounts

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
		SELECT "uid", "firstname", "email", "unsubscribe_token", "hash"
		FROM "user"
		WHERE
			"uid" = $1
	`, uid).Scan(&su.Uid, &su.Firstname, &su.Email, &su.UnsubToken, &hash); err != nil && err != sql.ErrNoRows {
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
		SELECT "uid", "firstname", "email", "unsubscribe_token", "hash"
		FROM "user"
		WHERE
			"email" = $1
	`, email).Scan(&su.Uid, &su.Firstname, &su.Email, &su.UnsubToken, &hash); err != nil && err != sql.ErrNoRows {
		return su, "", err
	}

	return su, hash, nil
}

// Create inserts the given account in database.
func (d *AccountDAO) Create(uid uuid.UUID, firstname, email, hash, unsubTok string, t time.Time) error {
	if _, err := d.DB.Exec(`
		INSERT INTO "user"
		("uid", "email", "firstname", "hash", "unsubscribe_token", "creation_time", "last_update")
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`, uid, email, firstname, hash, unsubTok, t, t); err != nil {
		return log.Err("Account/Create", err)
	}

	return nil
}

// Unsubscribe inserts an entry in the unsubscribe table indicating
// that an user doesn't wan't to receive any email anymore.
func (d *AccountDAO) Unsubscribe(token, reason string) error {
	// find the user having this unsubscribe token
	var err error
	var uid uuid.UUID

	if err = d.DB.QueryRow(`
		SELECT "uid" FROM "user"
		WHERE "unsubscribe_token" = $1
		LIMIT 1`, token).Scan(&uid); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Unsubscribe: unknown token: %s", token)
		} else {
			return err
		}
	}

	if uid.IsNil() {
		return fmt.Errorf("Unsubscribe: nil uid.")
	}

	if _, err = d.DB.Exec(`
		INSERT INTO "emailing_unsubscribe"
		(owner_uid, token, reason, creation_time)
		VALUES
		($1, $2, $3, $4)
	`, uid, token, reason, time.Now()); err != nil {
		return log.Err("Account/Unsubscribe", err)
	}

	return nil
}

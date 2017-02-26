// Account DAO.
//
// Rémy Mathieu © 2016

package accounts

import (
	"database/sql"
	"fmt"
	"time"

	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
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
		return nil, log.Err("UidByEmail", err)
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
		SELECT "uid", "firstname", "email", "unsubscribe_token", "stripe_token", "hash"
		FROM "user"
		WHERE
			"uid" = $1
	`, uid).Scan(&su.Uid, &su.Firstname, &su.Email, &su.UnsubToken, &su.StripeToken, &hash); err != nil && err != sql.ErrNoRows {
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
		SELECT "uid", "firstname", "email", "unsubscribe_token", "stripe_token", "hash"
		FROM "user"
		WHERE
			"email" = $1
	`, email).Scan(&su.Uid, &su.Firstname, &su.Email, &su.UnsubToken, &su.StripeToken, &hash); err != nil && err != sql.ErrNoRows {
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

// UpdateTz updates the timezone of the given user
// only if different.
func (d *AccountDAO) UpdateTz(uid uuid.UUID, tz string) error {
	_, err := d.DB.Exec(`
		UPDATE "user"
		SET
			"timezone" = $1,
			"last_update" = now()
		WHERE
			"uid" = $2
			AND
			"timezone" != $1
	`, tz, uid)

	if err != nil {
		return err
	}

	return nil
}

// UpdateStripeToken updates the Stripe token in database
// for the given user.
func (d *AccountDAO) UpdateStripeToken(u SimpleUser) error {
	n, err := d.DB.Exec(`
		UPDATE "user"
		SET
			"stripe_token" = $1,
			"last_update" = now()
		WHERE
			"uid" = $2
	`, u.StripeToken, u.Uid)

	if err != nil {
		return err
	}

	if n, err := n.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return fmt.Errorf("accounts: UpdateStripeToken: %d users updated.", n)
	}

	return nil
}

// UpdatePwdResetToken updates the password reset token in database
// for the given user. If the current token set in DB is still valid,
// we do not update the token and the time and return false.
func (d *AccountDAO) UpdatePwdResetToken(owner uuid.UUID, tok string, validUntil time.Time) (bool, error) {
	n, err := d.DB.Exec(`
		UPDATE "user"
		SET
			"password_reset_token" = $1,
			"password_reset_valid_until" = $2,
			"last_update" = now()
		WHERE
			"uid" = $3
			AND
			coalesce("password_reset_valid_until", '1970-01-01') < now()
	`, tok, validUntil, owner)

	if err != nil {
		return false, log.Err("UpdatePwdResetToken", err)
	}

	if n, err := n.RowsAffected(); err != nil {
		return false, log.Err("UpdatePwdResetToken", err)
	} else if n == 0 {
		return false, nil
	} else if n > 1 {
		return false, fmt.Errorf("accounts: UpdatePwdResetToken: %d users updated.", n)
	}

	return true, nil
}

// PwdReset updates the password and resets pwd reset fields
// for the given user. If the token is not valid anymore or
// the token is invalid, returns false.
func (d *AccountDAO) PwdReset(tok, pwd string) (bool, error) {
	n, err := d.DB.Exec(`
		UPDATE "user"
		SET
			"hash" = $1,
			"password_reset_token" = NULL,
			"password_reset_valid_until" = NULL,
			"last_update" = now()
		WHERE
			"password_reset_token" = $2
			AND
			"password_reset_valid_until" > now()
	`, pwd, tok)

	if err != nil {
		return false, log.Err("PwdReset", err)
	}

	if n, err := n.RowsAffected(); err != nil {
		return false, log.Err("PwdReset", err)
	} else if n == 0 {
		return false, nil
	} else if n > 1 {
		return false, fmt.Errorf("accounts: PwdReset: %d users updated.", n)
	}

	return true, nil
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

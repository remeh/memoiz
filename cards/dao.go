// Cards DAO.
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

// Update updates the text of the given
// card. It also updates the last_update time.
func (d *CardsDAO) UpdateText(uid, owner uuid.UUID, text string, t time.Time) (SimpleCard, error) {
	var position int

	if err := d.DB.QueryRow(`
		UPDATE "card"
		SET
			"text" = $1, "last_update" = $2
		WHERE
			"uid" = $3 AND "owner_uid" = $4
		RETURNING "position"
	`, text, t, uid, owner).Scan(&position); err != nil {
		return SimpleCard{}, err
	}

	return SimpleCard{
		Uid:      uid,
		Text:     text,
		Position: position,
	}, nil
}

// GetByUser returns the cards of the given user.
func (d *CardsDAO) GetByUser(uid uuid.UUID, state CardState) ([]SimpleCard, error) {
	rv := make([]SimpleCard, 0)

	rows, err := d.DB.Query(`
		SELECT "uid", "text", "position"
		FROM "card"
		WHERE
			"owner_uid" = $1
			AND
			"state" = $2
		ORDER BY "position" DESC
	`, uid, state)

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

// New creates a new card for the given user
// and returns its ID + position.
func (d *CardsDAO) New(userUid uuid.UUID, text string, t time.Time) (SimpleCard, error) {
	var rv SimpleCard

	uid := uuid.New()

	if err := d.DB.QueryRow(`
		INSERT INTO "card"
		("uid", "owner_uid", "text", "position", "creation_time", "last_update")
		SELECT $1, $2, $3, coalesce(max("position"),0)+1, $4, $4
		FROM "card"
		WHERE "owner_uid" = $2
		RETURNING "position"
	`, uid, userUid, text, t).Scan(&rv.Position); err != nil {
		return rv, fmt.Errorf("cards.New: %v", err)
	}

	rv.Uid = uid
	rv.Text = text

	return rv, nil
}

// Delete sets the deletion time of the given card in database
// and changes the state of the card.
func (d *CardsDAO) Delete(uid, owner uuid.UUID, t time.Time) error {
	if _, err := d.DB.Exec(`
		UPDATE "card"
		SET
			"deletion_time" = $1
		WHERE
			"uid" = $2
			AND
			"owner_uid" = $3
	`, t, uid, owner); err != nil {
		return fmt.Errorf("cards.Delete: %v", err)
	}
	return nil
}

func (d *CardsDAO) SwitchPosition(left, right, owner uuid.UUID, t time.Time) error {
	var tx *sql.Tx
	var err error

	if tx, err = d.DB.Begin(); err != nil {
		return fmt.Errorf("cards.SwitchPosition: can't start transaction: %v", err)
	}

	// retrieve actual position
	// ----------------------

	var lp, rp int64

	if err = tx.QueryRow(`
		SELECT "position" FROM "card"
		WHERE "uid" = $1 AND "owner_uid" = $2 FOR UPDATE`,
		left, owner).Scan(&lp); err != nil {
		tx.Rollback()
		return fmt.Errorf("cards.SwitchPosition: can't retrieve left card pos: %v", err)
	}

	if err = tx.QueryRow(`
		SELECT "position" FROM "card"
		WHERE "uid" = $1 AND "owner_uid" = $2 FOR UPDATE`,
		right, owner).Scan(&rp); err != nil {
		tx.Rollback()
		return fmt.Errorf("cards.SwitchPosition: can't retrieve right card pos: %v", err)
	}

	// set new position
	// ----------------------

	if _, err := tx.Exec(`
		UPDATE "card" SET "position" = $1
		WHERE "uid" = $2 AND "owner_uid" = $3`,
		rp, left, owner); err != nil {
		tx.Rollback()
		return fmt.Errorf("cards.SwitchPosition: can't update left card pos: %v", err)
	}

	if _, err := tx.Exec(`
		UPDATE "card" SET "position" = $1
		WHERE "uid" = $2 AND "owner_uid" = $3`,
		lp, right, owner); err != nil {
		tx.Rollback()
		return fmt.Errorf("cards.SwitchPosition: can't update right card pos: %v", err)
	}

	// commit the transaction
	// ----------------------

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("cards.SwitchPosition: can't commit transaction: %v", err)
	}

	return nil
}

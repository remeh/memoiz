// Memos DAO.
//
// Rémy Mathieu © 2016

package memos

import (
	"database/sql"
	"fmt"
	"time"

	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

type MemosDAO struct {
	DB *sql.DB
}

// ----------------------

var dao *MemosDAO

func DAO() *MemosDAO {
	if dao != nil {
		return dao
	}

	dao = &MemosDAO{
		DB: storage.DB(),
	}

	if err := dao.InitStmt(); err != nil {
		log.Error("Can't prepare MemosDAO")
		panic(err)
	}

	return dao
}

func (d *MemosDAO) InitStmt() error {
	var err error
	return err
}

// Update updates the text of the given
// memo. It also updates the last_update time.
func (d *MemosDAO) UpdateText(owner, uid uuid.UUID, text string, t time.Time) (Memo, error) {
	var position int

	if err := d.DB.QueryRow(`
		UPDATE "memo"
		SET
			"text" = $1, "last_update" = $2
		WHERE
			"uid" = $3 AND "owner_uid" = $4
		RETURNING "position"
	`, text, t, uid, owner).Scan(&position); err != nil {
		return Memo{}, log.Err("UpdateText:", err)
	}

	return Memo{
		Uid:      uid,
		Text:     text,
		Position: position,
	}, nil
}

// GetRichInfo returns the rich information added
// to the memo: category, link enrichment (img), etc.
func (d *MemosDAO) GetRichInfo(owner, uid uuid.UUID) (MemoRichInfo, error) {
	var ri MemoRichInfo

	if err := d.DB.QueryRow(`
		SELECT "r_category", "r_image", "r_url", "r_title", "last_update"
		FROM "memo"
		WHERE
			"uid" = $1
			AND
			"owner_uid" = $2
		`, uid, owner).Scan(&ri.Category, &ri.Image, &ri.Url, &ri.Title, &ri.LastUpdate); err != nil {
		if err == sql.ErrNoRows {
			return ri, nil
		}
		return ri, log.Err("GetRichInfo:", err)
	}

	return ri, nil
}

// Restore changes back the state of the memo to MemoActive
// and delete its archive time.
func (d *MemosDAO) Restore(owner, uid uuid.UUID, t time.Time) error {
	if _, err := d.DB.Exec(`
		UPDATE "memo"
		SET
			"archive_time" = $1,
			"state" = $2
		WHERE
			"uid" = $3
			AND
			"owner_uid" = $4
	`, nil, MemoActive, uid, owner); err != nil {
		return log.Err("memos.Restore", err)
	}
	return nil
}

// Archive changes the state of the memo to MemoArchived
// and sets its archive time.
func (d *MemosDAO) Archive(owner, uid uuid.UUID, t time.Time) error {
	if _, err := d.DB.Exec(`
		UPDATE "memo"
		SET
			"archive_time" = $1,
			"state" = $2
		WHERE
			"uid" = $3
			AND
			"owner_uid" = $4
	`, t, MemoArchived, uid, owner); err != nil {
		return log.Err("memos.Archive", err)
	}
	return nil
}

// UpdateLastNotif updates the last notification time of
// each given memo.
func (d *MemosDAO) UpdateLastNotif(uid uuid.UUID, memoUids []uuid.UUID, t time.Time) error {
	switch {
	case uid.IsNil():
		return fmt.Errorf("UpdateLastNotif: uid.IsNil()")
	case len(memoUids) == 0:
		return fmt.Errorf("UpdateLastNotif: len(memoUids) == 0")
	}

	vals := storage.Values(t, uid)
	for _, m := range memoUids {
		vals = append(vals, m)
	}

	if _, err := d.DB.Exec(`
		UPDATE "memo"
		SET
			"last_notification_time" = $1
		WHERE
			"owner_uid" = $2
			AND
			"uid" IN `+storage.InClause(len(memoUids))+`
	`, vals...); err != nil {
		return err
	}

	return nil
}

// GetByUser returns the memos of the given user.
func (d *MemosDAO) GetByUser(uid uuid.UUID, state MemoState) ([]Memo, error) {
	rv := make([]Memo, 0)

	rows, err := d.DB.Query(`
		SELECT "uid", "text", "position", "r_category", "r_image", "r_url", "r_title", "last_update"
		FROM "memo"
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
		var sc Memo
		var ri MemoRichInfo

		if err := rows.Scan(&sc.Uid, &sc.Text, &sc.Position, &ri.Category, &ri.Image, &ri.Url, &ri.Title, &ri.LastUpdate); err != nil {
			return rv, err
		}

		sc.MemoRichInfo = ri

		rv = append(rv, sc)
	}

	return rv, nil
}

// New creates a new memo for the given user
// and returns its ID + position.
func (d *MemosDAO) New(owner uuid.UUID, text string, t time.Time) (Memo, error) {
	var rv Memo

	memoUid := uuid.New()

	if err := d.DB.QueryRow(`
		INSERT INTO "memo"
		("uid", "owner_uid", "text", "position", "creation_time", "last_update")
		SELECT $1, $2, $3, coalesce(max("position"),0)+1, $4, $4
		FROM "memo"
		WHERE "owner_uid" = $2
		RETURNING "position"
	`, memoUid, owner, text, t).Scan(&rv.Position); err != nil {
		if err == sql.ErrNoRows {
			return rv, fmt.Errorf("memos.New: no position returned")
		}
		return rv, fmt.Errorf("memos.New: %v", err)
	}

	rv.Uid = memoUid
	rv.Text = text

	return rv, nil
}

// Delete sets the deletion time of the given memo in database
// and changes the state of the memo.
func (d *MemosDAO) Delete(owner, uid uuid.UUID, t time.Time) error {
	if _, err := d.DB.Exec(`
		UPDATE "memo"
		SET
			"deletion_time" = $1
		WHERE
			"uid" = $2
			AND
			"owner_uid" = $3
	`, t, uid, owner); err != nil {
		return fmt.Errorf("memos.Delete: %v", err)
	}
	return nil
}

func (d *MemosDAO) SwitchPosition(left, right, owner uuid.UUID, t time.Time) error {
	var tx *sql.Tx
	var err error

	if tx, err = d.DB.Begin(); err != nil {
		return fmt.Errorf("memos.SwitchPosition: can't start transaction: %v", err)
	}

	// retrieve actual position
	// ----------------------

	var lp, rp int64

	if err = tx.QueryRow(`
		SELECT "position" FROM "memo"
		WHERE "uid" = $1 AND "owner_uid" = $2 FOR UPDATE`,
		left, owner).Scan(&lp); err != nil {
		tx.Rollback()
		return fmt.Errorf("memos.SwitchPosition: can't retrieve left memo pos: %v", err)
	}

	if err = tx.QueryRow(`
		SELECT "position" FROM "memo"
		WHERE "uid" = $1 AND "owner_uid" = $2 FOR UPDATE`,
		right, owner).Scan(&rp); err != nil {
		tx.Rollback()
		return fmt.Errorf("memos.SwitchPosition: can't retrieve right memo pos: %v", err)
	}

	// set new position
	// ----------------------

	if _, err := tx.Exec(`
		UPDATE "memo" SET "position" = $1
		WHERE "uid" = $2 AND "owner_uid" = $3`,
		rp, left, owner); err != nil {
		tx.Rollback()
		return fmt.Errorf("memos.SwitchPosition: can't update left memo pos: %v", err)
	}

	if _, err := tx.Exec(`
		UPDATE "memo" SET "position" = $1
		WHERE "uid" = $2 AND "owner_uid" = $3`,
		lp, right, owner); err != nil {
		tx.Rollback()
		return fmt.Errorf("memos.SwitchPosition: can't update right memo pos: %v", err)
	}

	// commit the transaction
	// ----------------------

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("memos.SwitchPosition: can't commit transaction: %v", err)
	}

	return nil
}

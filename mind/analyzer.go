package mind

import (
	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"
)

type Analyzer interface {
	Fetch(string) error
	Analyze() (Categories, error)
	Store() error
}

func Analyze(uid uuid.UUID, text string) {
	if uuid.IsNil(uid) || len(text) == 0 {
		return
	}

	var a Analyzer
	var err error
	var cats Categories

	// TODO(remy): try from the cache "domain_result" first
	// TODO(remy): if the text is too long, should not be
	// useful to call Google Knowledge Graph
	a = &Kg{}

	if err = a.Fetch(text); err != nil {
		log.Error("Analyze/Fetch:", err)
		return
	}

	if cats, err = a.Analyze(); err != nil {
		log.Error("Analyze/Analyze:", err)
		return
	}

	if err = a.Store(); err != nil {
		log.Error("Analyze/Store:", err)
		return
	}

	if len(cats) == 0 || cats[0] == Unknown {
		return
	}

	// update the card
	if _, err := storage.DB().Exec(`
		UPDATE "card"
		SET "category" = $1
		WHERE "uid" = $2
	`, cats[0], uid); err != nil {
		log.Error("mind.Analyze:", err)
	}
}

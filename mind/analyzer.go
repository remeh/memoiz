package mind

import (
	"strings"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"
)

type Analyzer interface {
	TryCache(string) (bool, error)
	Fetch(string) error
	Analyze() error
	Store() error
	Categories() Categories
}

func Analyze(uid uuid.UUID, text string) {
	if uuid.IsNil(uid) || len(text) == 0 {
		return
	}

	var a Analyzer
	var err error

	// choose the first analyzer to launch
	// ----------------------

	spaces := strings.Count(text, " ")

	a = &Kg{}

	// do not use Google KG if it has too many spaces
	if spaces > 4 {
		a = &Bing{}
	}

	// don't even bother to analyze something which looks
	// like a complete not
	if spaces > 10 {
		a = &Stub{}
	}

	// apply the analyze
	// ----------------------

	found, err := a.TryCache(text)
	if err != nil {
		log.Error("Analyze/TryCache:", err)
		return
	}

	if !found {
		if err = a.Fetch(text); err != nil {
			log.Error("Analyze/Fetch:", err)
			return
		}
	}

	if err = a.Analyze(); err != nil {
		log.Error("Analyze/Analyze:", err)
		return
	}

	if err = a.Store(); err != nil {
		log.Error("Analyze/Store:", err)
		return
	}

	// update the Card if anything has been found
	// ----------------------

	cats := a.Categories()
	if cats == nil || len(cats) == 0 || cats[0] == Unknown {
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

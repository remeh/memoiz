package mind

import (
	"strings"

	"remy.io/scratche/config"
	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"
)

var (
	UseBing = false
	UseKg   = false
)

func init() {
	if len(config.Config.BingApiKey) > 0 {
		UseBing = true
		log.Info("Bing Search API will be used.")
	}

	if len(config.Config.KgApiKey) > 0 {
		UseKg = true
		log.Info("Google Knowledge Graph will be used.")
	}

}

// ----------------------

type Analyzer interface {
	TryCache(string) (bool, error)
	Fetch(string) error
	Analyze() error
	Store(cardUid uuid.UUID) error
	Categories() Categories
}

func Analyze(uid uuid.UUID, text string) {
	if uid.IsNil() || len(text) == 0 {
		return
	}

	analyzers := make([]Analyzer, 0)

	// choose the first analyzer to launch
	// ----------------------

	spaces := strings.Count(text, " ")

	// do not use Google KG if it has too many spaces
	if spaces <= 4 {
		analyzers = append(analyzers, &Kg{})
	}

	if spaces < 15 {
		analyzers = append(analyzers, &Bing{})
	}

	if len(analyzers) == 0 {
		// don't even bother to analyze something which looks
		// like a complete note
		analyzers = append(analyzers, &Stub{})
	}

	// always try to retrieve urls
	analyzers = append(analyzers, &Url{})

	// apply the analyze
	// ----------------------

	for _, a := range analyzers {
		var err error
		var found bool

		found, err = a.TryCache(text)
		if err != nil {
			log.Error("Analyze/TryCache:", err)
			return
		}

		if !found {
			if err = a.Fetch(text); err != nil {
				log.Error("Analyze/Fetch:", err)
				continue
			}
		}

		if err = a.Analyze(); err != nil {
			log.Error("Analyze/Analyze:", err)
			continue
		}

		if err = a.Store(uid); err != nil {
			log.Error("Analyze/Store:", err)
			continue
		}

		// update the Card if anything has been found
		// ----------------------

		cats := a.Categories()
		if cats != nil && len(cats) > 0 && cats[0] != Unknown {
			// update the card
			if _, err := storage.DB().Exec(`
			UPDATE "card"
			SET "r_category" = $1
			WHERE "uid" = $2
		`, cats[0], uid); err != nil {
				log.Error("mind.Analyze:", err)
			}

			return // we don't need to launch another analyzer
		}
	}
}

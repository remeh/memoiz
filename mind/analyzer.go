package mind

import (
	"strings"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
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
	Store(memoUid uuid.UUID) error
	Categories() Categories
}

// Analyze starts the regular analyzers such
// as Bing, Google Knowledge Graph.
func Analyze(uid uuid.UUID, text string) {
	if uid.IsNil() || len(text) == 0 {
		return
	}

	analyzers := make([]Analyzer, 0)

	// looks whether the text contains an URL
	// ---------------------

	url := rxUrl.FindString(text)
	if len(url) != 0 {
		// we have an URL: start only the URL
		// analyzer.
		analyzers = append(analyzers, &Url{url: url})
		analyze(analyzers, uid, text)
		return
	}

	// no URL, we'll try with Bing and Kg analyzers.
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
		// like a complete note.
		analyzers = append(analyzers, &Stub{})
	}

	analyze(analyzers, uid, text)
}

// ----------------------

// analyzer pipeline.
func analyze(analyzers []Analyzer, uid uuid.UUID, text string) {
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

		// update the Memo if anything has been found
		// ----------------------

		cats := a.Categories()
		if cats != nil && len(cats) > 0 && cats[0] != Uncategorized {
			// update the memo
			if _, err := storage.DB().Exec(`
			UPDATE "memo"
			SET "r_category" = $1
			WHERE "uid" = $2
		`, cats[0], uid); err != nil {
				log.Error("mind.Analyze:", err)
			}

			return // we don't need to launch another analyzer
		}
	}
}

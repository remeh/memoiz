package mind

import (
	"database/sql"
	"fmt"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
)

// guessByDomains retrieve the Category which seems to represent
// the best the given card.
// The weight the choice is also returned.
func guessByDomains(domains []string) (Category, int, error) {
	if len(domains) == 0 {
		log.Warning("guessByDomains: called with len(domains) == 0")
		return Unknown, 0, nil
	}

	inClause := "("
	for i := 0; i < len(domains); i++ {
		inClause += fmt.Sprintf("$%d", i+1)
		if i != len(domains)-1 {
			inClause += ","
		}
	}
	inClause += ")"

	var params []interface{} = make([]interface{}, len(domains))
	for i := range params {
		params[i] = domains[i]
	}

	var cat Category
	var weight int

	if err := storage.DB().QueryRow(fmt.Sprintf(`
		SELECT "category", sum("weight") w
		FROM "domain"
		WHERE "domain" IN
		%s
		GROUP BY "category"
		ORDER BY w
		DESC
		LIMIT 1
		`, inClause), params...).Scan(&cat, &weight); err != nil {
		if err != sql.ErrNoRows {
			err = fmt.Errorf("can't categorize: %v : %v", domains, err)
		}
	}

	// We do not want to return the category if there
	// is many domains tested and the resulting weight
	// is < 150
	if weight < 150 && len(domains) > 1 {
		return Unknown, 0, nil
	}

	return cat, weight, nil
}

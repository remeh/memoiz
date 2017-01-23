package mind

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/lib/pq"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"
)

const (
	KgUrl = "https://kgsearch.googleapis.com/v1/entities:search" // &indent=True
)

// Kg uses Google Knowledge Graph to give a Category
// to the given text.
type Kg struct {
	text        string
	types       []string
	description string
	categories  Categories
}

func (k *Kg) TryCache(text string) (bool, error) {
	if !UseKg {
		return false, nil
	}

	// TODO(remy): not implemented
	k.categories = Categories{Unknown}
	return false, nil
}

func (k *Kg) Fetch(text string) error {
	if !UseKg {
		return nil
	}

	var req *http.Request
	var resp *http.Response
	var err error

	k.text = text
	k.categories = Categories{Unknown}

	// http request to Bing
	// ----------------------

	if req, err = http.NewRequest("GET", k.buildUrl(), nil); err != nil {
		return err
	}

	cli := &http.Client{
		Timeout: time.Second * 15,
	}
	if resp, err = cli.Do(req); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http error %d when calling Kg", resp.StatusCode)
	}

	// read the response
	// ----------------------

	var data []byte
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	// read the types and description
	// ----------------------

	var elements []byte

	if elements, _, _, err = jsonparser.Get(data, "itemListElement"); err != nil {
		return err
	}

	jsonparser.ArrayEach(elements, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// description
		desc, _, _, _ := jsonparser.Get(value, "result", "description")
		k.description = strings.ToLower(string(desc))

		// types
		result, _, _, _ := jsonparser.Get(value, "result")
		t, _, _, _ := jsonparser.Get(result, "@type")
		jsonparser.ArrayEach(t, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			// ignore "Thing" too vague
			str := string(value)
			if str != "Thing" {
				k.types = append(k.types, strings.ToLower(string(value)))
			}
		})
	})

	return nil
}

func (k *Kg) Analyze() error {
	if !UseKg {
		return nil
	}

	if len(k.types) == 0 {
		return nil
	}

	// TODO(remy): use the description
	var c Category

	if err := storage.DB().QueryRow(`
		SELECT "category"
		FROM "kg_type"
		WHERE
			"type" = $1
	`, k.types[0]).Scan(&c); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		c = Unknown
	}

	k.categories = Categories{c}
	return nil
}

func (k *Kg) Store(memoUid uuid.UUID) error {
	if !UseKg {
		return nil
	}

	uid := uuid.New()

	// store
	if _, err := storage.DB().Exec(`
		INSERT INTO "kg_result"
		("uid", "memo_uid", "memo_text", "types", "description", "category", "creation_time")
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`, uid, memoUid, k.text, pq.Array(k.types), k.description, pq.Array(k.categories), time.Now()); err != nil {
		return err
	}

	// some log
	log.Debug("Kg decided that '", k.text, "' is '", k.categories)

	return nil
}

func (k *Kg) Categories() Categories {
	return k.categories
}

// ----------------------

func (k *Kg) buildUrl() string {
	return fmt.Sprintf("%s?limit=1&query=%s&key=%s", KgUrl, url.QueryEscape(k.text), config.Config.KgApiKey)
}

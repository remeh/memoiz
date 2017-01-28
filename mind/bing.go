package mind

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"remy.io/memoiz/config"
	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"

	"github.com/buger/jsonparser"
	"github.com/lib/pq"
)

const (
	BingUrl = "https://api.cognitive.microsoft.com/bing/v5.0/search?mkt=en-US"
)

// Bing uses Cognitive Services to do a Web Search
// and compares the returned domains to known ones
// to guess what could be the topic of the memo.
type Bing struct {
	text       string
	domains    []string
	categories Categories
	weight     int
}

func (b *Bing) TryCache(text string) (bool, error) {
	if !UseBing {
		return false, nil
	}

	// TODO(remy): not implemented
	b.categories = Categories{Uncategorized}
	return false, nil
}

func (b *Bing) Fetch(text string) error {
	if !UseBing {
		return nil
	}

	if len(text) < 2 {
		// nothing useful can be found with only two chars
		b.domains = make([]string, 0)
		return nil
	}

	var req *http.Request
	var resp *http.Response
	var err error

	b.text = text

	// http request to Bing
	// ----------------------

	if req, err = http.NewRequest("GET", b.buildUrl(), nil); err != nil {
		return err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", config.Config.BingApiKey)

	cli := &http.Client{}
	if resp, err = cli.Do(req); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http error %d when calling Bing", resp.StatusCode)
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

	// read the domains
	// ----------------------

	var webPages []byte

	if webPages, _, _, err = jsonparser.Get(data, "webPages", "value"); err != nil {
		return err
	}

	jsonparser.ArrayEach(webPages, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		d, _, _, _ := jsonparser.Get(value, "displayUrl")
		if url, err := url.Parse(strings.Replace(string(d), `\/`, "/", -1)); err == nil {
			if len(url.Host) == 0 {
				b.domains = append(b.domains, b.extractDomain(url.Path))
			} else {
				b.domains = append(b.domains, b.extractDomain(url.Host))
			}
		}
	})

	return nil
}

func (b *Bing) Analyze() error {
	if !UseBing {
		return nil
	}

	cat, weight, err := guessByDomains(b.domains)
	if err != nil {
		log.Debug("Bing.Analyze:", err)
	}
	b.weight = weight
	b.categories = Categories{cat}

	return nil
}

func (b *Bing) Store(memoUid uuid.UUID) error {
	if !UseBing {
		return nil
	}

	uid := uuid.New()

	// store
	if _, err := storage.DB().Exec(`
		INSERT INTO "domain_result"
		("uid", "memo_uid", "memo_text", "category", "domains", "weight", "creation_time")
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`, uid, memoUid, b.text, pq.Array(b.categories), pq.Array(b.domains), b.weight, time.Now()); err != nil {
		return err
	}

	// some log
	log.Debug("Bing decided that '", b.text, "' is '", b.categories, "' (weight:", b.weight, ")")

	return nil
}

func (b *Bing) Categories() Categories {
	return b.categories
}

// ----------------------

// buildUrl returns the URL to call Bing API.
func (b *Bing) buildUrl() string {
	return fmt.Sprintf("%s&count=30&q=%s", BingUrl, url.QueryEscape(b.text))
}

// extractDomain extracts the domain from the given URL.
func (b *Bing) extractDomain(url string) string {
	// NOTE(remy): we add the / for the regexp
	// NOTE(remy): we take only the first match, this is why
	// I don't use FindAllStringSubmatch
	str := rxDomain.FindStringSubmatch(strings.ToLower(url) + "/")
	if len(str) == 2 {
		return str[1]
	}
	return ""
}

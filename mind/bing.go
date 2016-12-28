package mind

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"remy.io/scratche/config"
	"remy.io/scratche/log"
	"remy.io/scratche/storage"

	"github.com/buger/jsonparser"
)

const (
	BingUrl = "https://api.cognitive.microsoft.com/bing/v5.0/search?mkt=en-US"
)

// Bing uses Cognitive Services to do a Web Search
// and compares the returned domains to known ones
// to guess what could be the topic of the card.
type Bing struct {
	text       string
	domains    []string
	categories Categories
}

func (b *Bing) Fetch(text string) error {
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

func (b *Bing) Analyze() (Categories, error) {
	cat, err := b.guessByDomains()
	if err != nil {
		log.Debug("domains:", b.domains)
		return Categories{Unknown}, fmt.Errorf("Bing.Analyze: %v", err)
	}

	b.categories = Categories{cat}

	return b.categories, err
}

func (b *Bing) Store() error {
	// TODO(remy): store the score for this guess
	log.Debug("Bing decided that '", b.text, "' is '", b.categories, "'")
	if b.categories[0] == Unknown {
		log.Debug("domains:", b.domains)
	}
	// not implemented
	return nil
}

// ----------------------

// rxDomain retrieves only the domain (removing the TLD)
var rxDomain *regexp.Regexp = regexp.MustCompile(`([a-zA-Z0-9]*)\.[a-zA-Z0-9]*\/`)

// guessByDomains retrieve the Category which seems to represent
// the best the given card.
func (b *Bing) guessByDomains() (Category, error) {
	if len(b.domains) == 0 {
		// TODO(remy): log this, shouldn't happen if bing has respond
		return Unknown, nil
	}

	inClause := "("
	for i := 0; i < len(b.domains); i++ {
		inClause += fmt.Sprintf("$%d", i+1)
		if i != len(b.domains)-1 {
			inClause += ","
		}
	}
	inClause += ")"

	fmt.Println(b.domains)

	var params []interface{} = make([]interface{}, len(b.domains))
	for i := range params {
		params[i] = b.domains[i]
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
		return Unknown, fmt.Errorf("can't categorize: %v : %v", b.domains, err)
	}

	if weight < 250 {
		return Unknown, nil
	}

	return cat, nil
}

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

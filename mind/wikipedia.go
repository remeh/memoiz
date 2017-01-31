package mind

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Wikipedia struct {
	text string
	cat  Category
	data []byte
}

var rxDelBrackets = regexp.MustCompile(`(\[.*?\]) *`)

func (w *Wikipedia) Enrich(text string, cat Category) (bool, EnrichResult, error) {

	w.text = text
	w.cat = cat

	if !w.validate() {
		return false, EnrichResult{}, nil
	}

	// TODO(remy): ensure there is no disambiguity with other pages

	// get the content

	if found, err := w.fetchContent(); !found || err != nil {
		return found, EnrichResult{}, err
	}

	return w.extract()
}

// ----------------------

// validate returns whether or not this text/category
// is applicable
func (w *Wikipedia) validate() bool {
	if strings.Count(w.text, " ") > 4 {
		return false
	}
	return true
}

func (w *Wikipedia) extract() (bool, EnrichResult, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(w.data))
	if err != nil {
		return false, EnrichResult{}, err
	}

	// read the meta
	// ----------------------

	var result EnrichResult

	if selection := doc.Find("#mw-content-text p").First(); selection != nil {
		// remove text inside []
		str := selection.Text()
		result.Content = string(rxDelBrackets.ReplaceAll([]byte(str), []byte{}))
	}

	if selection := doc.Find("#footer-info-copyright").First(); selection != nil {
		result.ContentCopyright = contentCopyright(selection.Text())
	}

	// TODO(remy): imgUrl and image license

	return len(result.Content) != 0, result, nil
}

func (w *Wikipedia) fetchContent() (bool, error) {

	// TODO(remy): try to fetch the page, could be not found!
	// ----------------------

	var err error
	var req *http.Request
	var resp *http.Response

	url := w.generateUrl()

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", randomUserAgent())

	cli := &http.Client{
		Timeout: time.Second * 15,
	}
	if resp, err = cli.Do(req); err != nil {
		return false, err
	}

	if resp.StatusCode == 404 {
		return false, nil
	} else if resp.StatusCode != 200 {
		return false, fmt.Errorf("http error %d when retrieving %s", resp.StatusCode, url)
	}

	// read the response
	// ----------------------

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if w.data, err = ioutil.ReadAll(resp.Body); err != nil {
		return false, err
	}

	return true, nil
}

func (w *Wikipedia) generateUrl() string {
	return "https://en.wikipedia.org/wiki/" + strings.Replace(url.QueryEscape(w.text), "+", "%20", -1)
}

func contentCopyright(str string) string {
	// we do not want this part of the response
	str = strings.Replace(str, " By using this site, you agree to the Terms of Use and Privacy Policy. Wikipedia® is a registered trademark of the Wikimedia Foundation, Inc., a non-profit organization.", "", -1)
	str = strings.Replace(str, "\n", " ", -1)
	return str
}

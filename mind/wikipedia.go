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

func (w *Wikipedia) Enrich(text string, cat Category) (bool, Description, ImageUrl, error) {

	w.text = text
	w.cat = cat

	if !w.validate() {
		return false, "", "", nil
	}

	// TODO(remy): ensure there is no disambiguity with other pages

	// get the content

	if found, err := w.fetchContent(); !found || err != nil {
		return found, "", "", err
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

func (w *Wikipedia) extract() (bool, Description, ImageUrl, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(w.data))
	if err != nil {
		return false, "", "", err
	}

	// read the meta
	// ----------------------

	var found bool
	var desc Description
	var imgUrl ImageUrl

	doc.Find("#mw-content-text p").Each(func(i int, s *goquery.Selection) {
		if i != 0 { // we only want the first one
			return

		}
		// remove text inside []
		str := s.Text()
		desc = Description(rxDelBrackets.ReplaceAll([]byte(str), []byte{}))

		// TODO(remy): imgUrl
		found = len(desc) != 0
	})

	return true, desc, imgUrl, nil
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

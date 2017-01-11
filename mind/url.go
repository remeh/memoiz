// TODO(remy): not sure this should be an analyzer

package mind

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"

	"remy.io/scratche/log"
	"remy.io/scratche/storage"
	"remy.io/scratche/uuid"

	"github.com/PuerkitoBio/goquery"
)

// E.g. a Tweet:
// <meta property="og:type" content="article">
// <meta property="og:url" content="https://twitter.com/troyhunt/status/800632175880183808">
// <meta property="og:title" content="Troy Hunt on Twitter">
// <meta property="og:image" content="https://pbs.twimg.com/profile_images/1154092629/Square__small__400x400.jpg">
// <meta property="og:description" content="“Just blogged: Ad blockers are part of the problem https://t.co/OrQXK7QxZ8”">
// <meta property="og:site_name" content="Twitter">
// <meta property="fb:app_id" content="2231777543">
const (
	metaUrl         = "og:url"
	metaImage       = "og:image"
	metaTitle       = "og:title"
	metaDescription = "og:description"
)

type Url struct {
	url   string
	image string
	title string
	data  []byte
}

var urlRx = regexp.MustCompile(`((https?:\/\/)?([0-9a-zA-Z]+\.)*[-_0-9a-zA-Z]+\.[-_0-9a-zA-Z]+)\/([-_0-9a-zA-Z\.\/])*(\?[0-9a-zA-Z\%\&\-\=\_\.]*)`)

func (u *Url) TryCache(text string) (bool, error) {
	return false, nil
}

func (u *Url) Fetch(text string) error {
	// look for an url in the text
	u.url = urlRx.FindString(text)
	if len(u.url) == 0 {
		return nil
	}

	log.Debug("Url is fetching", u.url)

	// url found, try to fetch the page
	// ----------------------

	var err error
	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest("GET", u.url, nil); err != nil {
		return err
	}
	req.Header.Set("User-Agent", randomUserAgent())

	cli := &http.Client{} // TODO(remy): parameters of this client ?
	// XXX(remy): add timeout to the client
	if resp, err = cli.Do(req); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http error %d when retrieving %s", resp.StatusCode, u.url)
	}

	// read the response
	// ----------------------

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if u.data, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	return nil
}

func (u *Url) Analyze() error {
	if u.data == nil {
		return nil
	}

	// read the fetch data
	// TODO(remy): we should probably ensure its html first ?
	// TODO(remy): pretty sure we don't need to read the whole file

	// read the document as HTML

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(u.data))
	if err != nil {
		return err
	}

	var title, ogTitle, ogDescription string

	// read the title
	// ----------------------

	title = doc.Find("title").Text()

	// read the meta
	// ----------------------

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop, exists := s.Attr("property")
		if exists {
			switch prop {
			case metaUrl:
				u.url, _ = s.Attr("content")
			case metaImage:
				u.image, _ = s.Attr("content")
			case metaTitle:
				ogTitle, _ = s.Attr("content")
			case metaDescription:
				ogDescription, _ = s.Attr("content")
			}
		}
	})

	u.title = chooseTitle(u.url, title, ogTitle, ogDescription)

	return nil
}

func (u *Url) Store(uid uuid.UUID) error {
	// update the card image and URL
	if _, err := storage.DB().Exec(`
		UPDATE "card"
		SET "r_url" = $1, "r_image" = $2
		WHERE "uid" = $3
	`, u.url, u.image, uid); err != nil {
		return err
	}
	return nil
}

func (u *Url) Categories() Categories {
	return Categories{Unknown}
}

// ----------------------

// TODO(remy): comment me
func chooseTitle(url string, title, ogTitle, ogDescription string) string {
	var rv string

	// TODO(remy): get domain only from url
	// TODO(remy): depending on the domain, return the ogTitle or the ogDescription

	// fallback on title
	if len(rv) == 0 {
		rv = title
	}

	return rv
}

var uas []string = []string{
	// Chrome
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.99 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36",
	// Firefox
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12;rv:49.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12;rv:49.0) Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12;rv:49.0) Gecko/20100101 Firefox/48.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12;rv:49.0) Gecko/20100101 Firefox/49.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12;rv:49.0) Gecko/20100101 Firefox/50.0",
	// Internet Explorer
	"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
}

func randomUserAgent() string {
	return uas[rand.Intn(len(uas))]
}

// TODO(remy): not sure this should be an analyzer

package mind

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"remy.io/memoiz/log"
	"remy.io/memoiz/storage"
	"remy.io/memoiz/uuid"

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
	url    string
	domain string
	image  string
	title  string // short desc
	desc   string // long desc
	data   []byte

	category Category
}

func (u *Url) TryCache(text string) (bool, error) {
	return false, nil
}

func (u *Url) Fetch(text string) error {
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

	// NOTE(remy): we force ipv4 because youtube answers 429 responses (too many requests)
	// on OVH network if we're using ipv6...
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.Dial("tcp4", addr)
		},
	}

	cli := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 15,
	}
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

	matches := rxDomain.FindAllStringSubmatch(u.url, 2)
	if len(matches) >= 1 && len(matches[0]) >= 2 {
		u.domain = strings.ToLower(matches[0][1])
	}

	// first, checks if we can find a category
	// for this url.
	// ----------------------
	var cat Category
	var weight int
	var err error

	if cat, weight, err = guessByDomains([]string{u.domain}); err != nil {
		log.Error("Url/Analyze:", err)
		// we do not return, we want to try to fetch the URL
	}

	if cat != Uncategorized && weight > 50 {
		u.category = cat
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

	u.title, u.desc = titleAndDesc(u.domain, title, ogTitle, ogDescription)

	return nil
}

func (u *Url) Store(uid uuid.UUID) error {
	// update the memo image and URL
	if _, err := storage.DB().Exec(`
		UPDATE "memo"
		SET "r_url" = $1, "r_image" = $2, "r_title" = $3
		WHERE "uid" = $4
	`, u.url, u.image, u.title, uid); err != nil {
		return log.Err("Url:", err)
	}
	if u.category != Uncategorized {
		if _, err := storage.DB().Exec(`
		UPDATE "memo"
		SET "r_category" = $1
		WHERE "uid" = $2
	`, u.category, uid); err != nil {
			return err
		}
	}
	return nil
}

func (u *Url) Categories() Categories {
	return Categories{Uncategorized}
}

func (u *Url) Enrich(text string, cat Category) (bool, EnrichResult, error) {
	if err := u.Fetch(text); err != nil {
		return false, EnrichResult{}, err
	}

	if err := u.Analyze(); err != nil {
		return false, EnrichResult{}, err
	}

	if !validImage(u.image) {
		u.image = ""
	}

	rv := EnrichResult{
		ImageUrl: u.image,
		Content:  u.desc,
		Title:    u.title,
		Format:   EnrichStandard,
	}

	if text == u.url && len(u.image) == 0 {
		rv.Format = EnrichUrlNoImage
	}

	switch u.domain {
	case "youtube", "vimeo":
		if len(u.image) > 0 && len(u.title) > 0 {
			rv.Format = EnrichUrlFocusImage
		}
	}

	if len(u.desc) == 0 { // when no desc, use the title as description instead
		rv.Content = rv.Title
		rv.Title = ""
	}

	return true, rv, nil
}

// ----------------------

// titleAndDesc chooses the best title for the given URL
// and given read title, og:title and og:description.
func titleAndDesc(domain string, title, ogTitle, ogDescription string) (string, string) {
	var t string
	var d string

	if len(domain) != 0 {
		switch domain {
		case "twitter":
			t = ogDescription
			d = ogTitle
		default:
			t = ogTitle
			d = ogDescription
		}
	}

	if len(t) == 0 && len(ogDescription) != 0 {
		t = ogDescription
		d = ""
	}

	// fallback on title
	if len(t) == 0 {
		t = title
		d = ""
	}

	return t, d
}

// validImage returns whether or not this image
// should be considered valid to send.
func validImage(imgUrl string) bool {
	for _, img := range igImages {
		if img == imgUrl {
			return false
		}
	}
	return true
}

// igImages are images to ignored, they're not valid
// to send to the user.
var igImages []string = []string{
	"https://s0.wp.com/i/blank.jpg",
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

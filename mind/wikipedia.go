package mind

import (
	"bytes"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
)

type Wikipedia struct {
	text string
	cat  Category

	// contentData is the whole wikipedia
	// page fetched.
	contentData []byte

	imagesUrls     []string
	imagesLicenses []string
}

var rxDelBrackets = regexp.MustCompile(`(\[.*?\]) *`)

var licenses = []string{
	"CC BY-SA 3.0",
}

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

	if found, err := w.fetchImages(); !found || err != nil {
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

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(w.contentData))
	if err != nil {
		return false, EnrichResult{}, err
	}

	// read the first paragraph
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

	// extract the image
	// ----------------------

	// TODO(remy): imgUrl and image license
	// choose the Url depending on license, size?

	return len(result.Content) != 0, result, nil
}

func (w *Wikipedia) fetchContent() (bool, error) {

	// TODO(remy): try to fetch the page, could be not found!
	// ----------------------

	var err error
	var resp *http.Response

	if resp, err = Fetch(w.generateContentUrl(w.text)); err != nil {
		return false, err
	} else if resp == nil {
		return false, nil
	}

	// read the response
	// ----------------------

	var data []byte

	if data, err = Read(resp); err != nil {
		return false, err
	}
	w.contentData = data

	return true, nil
}

func (w *Wikipedia) fetchImages() (bool, error) {

	// http call
	// ----------------------

	var err error
	var resp *http.Response

	if resp, err = Fetch(w.generateImagesUrl(w.text)); err != nil {
		return false, err
	} else if resp == nil {
		return false, nil
	}

	// read the response
	// ----------------------

	var data []byte
	var elements []byte

	if data, err = Read(resp); err != nil {
		return false, err
	}

	// TODO(remy): stantardize the page id
	if elements, _, _, err = jsonparser.Get(data, "query", "pages", "901022", "images"); err != nil {
		return false, err
	}

	// look for images titles.
	// ----------------------

	var titles []string

	jsonparser.ArrayEach(elements, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// ignore this file
		title, err := jsonparser.GetString(value, "title")
		if err != nil {
			return
		}

		if len(title) == 0 {
			return
		}

		if title == "File:Commons-logo.svg" { // ignore this image
			return
		}

		titles = append(titles, title)
	})

	// gets each image license
	// ----------------------

	if len(titles) == 0 {
		return false, nil
	}

	for _, image := range titles {
		// look for the field license
		// ----------------------

		if resp, err = Fetch(w.generateImageLicense(image)); err != nil {
			return false, err
		} else if resp == nil {
			return false, nil
		}

		var err error
		var data []byte

		if data, err = Read(resp); err != nil {
			return false, err
		}

		license, err := jsonparser.GetString(data, "query", "pages", "-1", "imageinfo", "[0]", "extmetadata", "LicenseShortName", "value")
		if err != nil {
			return false, err
		}

		// look for its url
		// ----------------------

		if resp, err = Fetch(w.generateImageUrl(image)); err != nil {
			return false, err
		} else if resp == nil {
			return false, nil
		}

		if data, err = Read(resp); err != nil {
			return false, err
		}

		url, err := jsonparser.GetString(data, "query", "pages", "-1", "imageinfo", "[0]", "url")
		if err != nil {
			return false, err
		}

		// ok!
		// ----------------------

		w.imagesLicenses = append(w.imagesLicenses, license)
		w.imagesUrls = append(w.imagesUrls, url)
	}

	return true, nil
}

func (w *Wikipedia) generateImageLicense(image string) string {
	return "https://en.wikipedia.org/w/api.php?action=query&format=json&prop=imageinfo&iiprop=extmetadata&titles=" + w.escape(image)
}

func (w *Wikipedia) generateImageUrl(image string) string {
	return "https://en.wikipedia.org/w/api.php?action=query&format=json&prop=imageinfo&iiprop=url&titles=" + w.escape(image)
}

func (w *Wikipedia) generateImagesUrl(title string) string {
	return "https://en.wikipedia.org/w/api.php?action=query&prop=images&format=json&titles=" + w.escape(title)
}

func (w *Wikipedia) generateContentUrl(title string) string {
	return "https://en.wikipedia.org/wiki/" + w.escape(title)
}

func (w *Wikipedia) escape(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// ----------------------

func contentCopyright(str string) string {
	// we do not want this part of the response
	str = strings.Replace(str, " By using this site, you agree to the Terms of Use and Privacy Policy. Wikipedia® is a registered trademark of the Wikimedia Foundation, Inc., a non-profit organization.", "", -1)
	str = strings.Replace(str, "\n", " ", -1)
	return str
}

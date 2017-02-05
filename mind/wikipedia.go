package mind

import (
	"bytes"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"remy.io/memoiz/log"

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
	imagesFiles    []string
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

	if err := w.redirect(); err != nil {
		return false, EnrichResult{}, err
	}

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

// redirects queries wikipedia to know whether we should
// fetch another title name or not.
func (w *Wikipedia) redirect() error {

	// TODO(remy): try to fetch the page, could be not found!
	// ----------------------

	var err error
	var resp *http.Response

	if resp, err = Fetch(w.generateRedirects(w.text)); err != nil {
		return err
	} else if resp == nil {
		return nil
	}

	// read the response
	// ----------------------

	var data []byte

	if data, err = Read(resp); err != nil {
		return err
	}

	var title string

	if title, err = jsonparser.GetString(data, "query", "redirects", "[0]", "to"); err != nil {
		// no redirection, it's not really an error
		// we silent this because we can't do better
		return nil
	}

	log.Debug("Wikipedia:", w.text, "redirected to", title)
	w.text = title
	return nil
}

func (w *Wikipedia) extract() (bool, EnrichResult, error) {

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(w.contentData))
	if err != nil {
		return false, EnrichResult{}, err
	}

	// read the first paragraph
	// ----------------------

	var result EnrichResult

	if selection := doc.Find("#mw-content-text > p").First(); selection != nil {
		// remove text inside []
		str := selection.Text()
		result.Content = string(rxDelBrackets.ReplaceAll([]byte(str), []byte(" ")))
	}

	if selection := doc.Find("#footer-info-copyright").First(); selection != nil {
		result.ContentCopyright = contentCopyright(selection.Text())
		result.ContentSource = w.generateContentUrl(w.text)
	}

	if len(result.Content) == 0 {
		return false, result, nil // do not continue if we've not found any content
	}

	// extract the image
	// ----------------------

	if len(w.imagesUrls) != len(w.imagesLicenses) ||
		len(w.imagesUrls) != len(w.imagesFiles) {
		log.Error("Wikipedia: extract: len(w.imagesUrls) != len(w.imagesLicenses) ignoring image")
		return true, result, nil
	}

	if len(w.imagesUrls) == 0 {
		return true, result, nil
	}

	for i, license := range w.imagesLicenses {
		if !acceptedLicense(license) {
			continue
		}

		// TODO(remy): try to not send epicly large images

		result.ImageUrl = w.imagesUrls[i]
		result.ImageCopyright = license // TODO(remy): put a copyright notice, not the license title
		result.ImageSource = w.generateContentUrl(w.imagesFiles[i])
	}

	return len(result.Content) != 0, result, nil
}

func (w *Wikipedia) fetchContent() (bool, error) {

	// TODO(remy): try to fetch the page, could not existing!
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

	if buf, _, _, err := jsonparser.Get(data, "query", "pages"); err != nil {
		return false, err
	} else {
		i := 0 // we will only take the first successful
		jsonparser.ObjectEach(buf, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if i != 0 {
				return nil
			}

			if d, _, _, err := jsonparser.Get(value, "images"); err != nil {
				return err
			} else {
				elements = d
			}

			i++
			return nil
		})
	}

	// look for images titles.
	// ----------------------

	var files []string

	jsonparser.ArrayEach(elements, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// ignore this file
		title, err := jsonparser.GetString(value, "title")
		if err != nil {
			return
		}

		if len(title) == 0 {
			return
		}

		if title == "File:Commons-logo.svg" || // ignore this image
			strings.HasSuffix(strings.ToLower(title), "svg") { // no support for SVG in gmail
			return
		}

		files = append(files, title)
	})

	// gets each image license
	// ----------------------

	if len(files) == 0 {
		return false, nil
	}

	for _, file := range files {
		// look for the field license
		// ----------------------

		if resp, err = Fetch(w.generateImageLicense(file)); err != nil {
			return false, err
		} else if resp == nil {
			return false, nil
		}

		var err error
		var data []byte
		var license string

		if data, err = Read(resp); err != nil {
			return false, err
		}

		if buf, _, _, err := jsonparser.Get(data, "query", "pages"); err != nil {
			return false, err
		} else {
			i := 0 // we will only take the first successful
			jsonparser.ObjectEach(buf, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				if i != 0 {
					return nil
				}

				if d, err := jsonparser.GetString(value, "imageinfo", "[0]", "extmetadata", "LicenseShortName", "value"); err != nil {
					return err
				} else {
					license = d
				}

				i++
				return nil
			})
		}

		// look for its url
		// ----------------------

		var url string

		if resp, err = Fetch(w.generateImageUrl(file)); err != nil {
			return false, err
		} else if resp == nil {
			return false, nil
		}

		if data, err = Read(resp); err != nil {
			return false, err
		}

		if buf, _, _, err := jsonparser.Get(data, "query", "pages"); err != nil {
			return false, err
		} else {
			i := 0 // we will only take the first successful
			jsonparser.ObjectEach(buf, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				if i != 0 {
					return nil
				}

				if d, err := jsonparser.GetString(value, "imageinfo", "[0]", "url"); err != nil {
					return err
				} else {
					url = d
				}

				i++
				return nil
			})
		}

		// ok!
		// ----------------------

		w.imagesLicenses = append(w.imagesLicenses, license)
		w.imagesUrls = append(w.imagesUrls, url)
		w.imagesFiles = append(w.imagesFiles, file) // e.g. File:Streetlight%20Manifesto.jpg
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

func (w *Wikipedia) generateRedirects(title string) string {
	return "https://en.wikipedia.org/w/api.php?action=query&redirects&format=json&titles=" + w.escape(title)
}

func (w *Wikipedia) escape(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// ----------------------

func acceptedLicense(license string) bool {
	for _, l := range licenses {
		if l == license {
			return true
		}
	}
	return false
}

func contentCopyright(str string) string {
	// we do not want this part of the response
	str = strings.Replace(str, " By using this site, you agree to the Terms of Use and Privacy Policy. Wikipedia® is a registered trademark of the Wikimedia Foundation, Inc., a non-profit organization.", "", -1)
	str = strings.Replace(str, "\n", " ", -1)
	return str
}

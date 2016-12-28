package mind

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/buger/jsonparser"

	"remy.io/scratche/config"
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
}

func (k *Kg) Fetch(text string) error {
	// TODO(remy): test text for amount of spaces

	var req *http.Request
	var resp *http.Response
	var err error

	k.text = text

	// http request to Bing
	// ----------------------

	if req, err = http.NewRequest("GET", k.buildUrl(), nil); err != nil {
		return err
	}

	fmt.Println(k.buildUrl())

	cli := &http.Client{}
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
			k.types = append(k.types, strings.ToLower(string(value)))
		})
	})

	return nil
}

func (k *Kg) Analyze() (Categories, error) {
	fmt.Println(k.description)
	fmt.Println(k.types)
	return Categories{Unknown}, nil
}

func (k *Kg) Store() error {
	return nil
}

func (k *Kg) buildUrl() string {
	return fmt.Sprintf("%s?limit=1&query=%s&key=%s", KgUrl, url.QueryEscape(k.text), config.Config.KgApiKey)
}

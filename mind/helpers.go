package mind

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// fetch fetches the given URL. Returns nil response on 404.
func Fetch(url string) (*http.Response, error) {
	var err error
	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", randomUserAgent())

	cli := &http.Client{
		Timeout: time.Second * 15,
	}
	if resp, err = cli.Do(req); err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error %d when retrieving %s", resp.StatusCode, url)
	}

	return resp, nil
}

// Read reads the content of an http response.
func Read(r *http.Response) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("mind: Read: called with a nil response")
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	return ioutil.ReadAll(r.Body)
}

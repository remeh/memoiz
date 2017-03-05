package mind

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var licenses = []string{
	"CC BY-SA 3.0",
	"http://creativecommons.org/licenses/by-sa/2.0",
	"http://creativecommons.org/licenses/by-sa/2.5",
	"http://creativecommons.org/licenses/by-sa/3.0",
	"http://creativecommons.org/licenses/by-sa/4.0",
	"http://creativecommons.org/licenses/by/2.0",
	"http://creativecommons.org/licenses/by/2.5",
	"http://creativecommons.org/licenses/by/3.0",
	"http://creativecommons.org/licenses/by/4.0",
	"https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License",
}

// ----------------------

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

// AcceptedLicense returns whether the given license
// valid for us (meaning we can use the content).
func AcceptedLicense(license string) bool {
	for _, l := range licenses {
		if l == license {
			return true
		}
	}
	return false
}

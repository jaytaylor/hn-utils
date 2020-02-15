package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	// BaseURL holds the base URL for HackerNews.
	BaseURL = "https://news.ycombinator.com"

	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36"
)

// Login returns an authenticated *http.Client (or errors out).
func Login(username string, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %s", err)
	}

	client := NoAuthClient()
	client.Jar = jar

	form := url.Values{}
	form.Add("acct", username)
	form.Add("pw", password)
	form.Add("goto", "news")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", BaseURL), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("login: creating POST request: %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", BaseURL+"/")
	req.Header.Add("Origin", BaseURL+"/")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Authority", strings.Split(BaseURL, "://")[1])
	req.Header.Add("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 3 || resp.Header.Get("Location") == "" {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("login: expected 3xx reponse status-code but got %v (body=%v)", resp.StatusCode, string(body))
	}
	return client, nil
}

func NoAuthClient() *http.Client {
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

// CheckedGet requires an already authenticated *http.Client and retrieves
// content from the specified page.
func CheckedGet(client *http.Client, page string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		return nil, fmt.Errorf("creating logged-in GET request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting %v: %s", page, err)
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("expected 2xx response status-code from %v but got %v", page, resp.StatusCode)
	}

	return resp.Body, nil
}

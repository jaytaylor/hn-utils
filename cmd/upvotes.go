package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/jaytaylor/hn-utils/domain"

	"gigawatt.io/ago"
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36"

var (
	Password string
)

func initUpvotesFlags() {
	upvotesCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "HN login password")
}

var upvotesCmd = &cobra.Command{
	Use:   "upvotes",
	Short: "HN user upvotes",
	Long:  "Retrieves upvotes for a given HN user/password",
	PreRun: func(_ *cobra.Command, _ []string) {
		initLogging()

		if Password == "" {
			errorExit(errors.New("Missing required flag: -p/--password must not be empty"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			stories  = []domain.Story{}
			moreLink = fmt.Sprintf("%v/upvoted?id=%v", baseURL, User)
		)

		client, err := login(User, Password)
		if err != nil {
			errorExit(err)
		}
		log.Debug("Logged in successfully")

		for len(moreLink) > 0 {
			now := time.Now()

			log.WithField("more-link", moreLink).Debug("Fetching")

			rc, err := getLoggedInPage(client, moreLink)
			if err != nil {
				errorExit(err)
			}

			doc, err := goquery.NewDocumentFromReader(rc)
			if err != nil {
				errorExit(err)
			}
			if err := rc.Close(); err != nil {
				errorExit(fmt.Errorf("closing response body from %v: %s", moreLink, err))
			}

			doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
				var (
					title    = s.Find(".title a.storylink")
					comments = s.Next().Find("a").Last()
					story    = domain.Story{
						ID:          int64Or(s.AttrOr("id", "0"), -1),
						Title:       title.Text(),
						URL:         reconstructHNURL(title.AttrOr("href", "")),
						Points:      int64Or(numExpr.ReplaceAllString(s.Next().Find(".score").Text(), "$1"), -1),
						Comments:    int64Or(numExpr.ReplaceAllString(comments.Text(), "$1"), -1),
						CommentsURL: comments.AttrOr("href", ""),
						Submitter:   hnuserExpr.ReplaceAllString(s.Next().Find(".hnuser").AttrOr("href", ""), ""),
					}
				)
				if len(story.CommentsURL) > 0 && !strings.HasPrefix(story.CommentsURL, "https://") {
					story.CommentsURL = fmt.Sprintf("%s/%s", baseURL, story.CommentsURL)
				}

				humanTime := s.Next().Find(".age").Text()
				// For items favorited in the early days of the feature, HN spits out
				// "on <date>" instead of a human-style delta.
				switch strings.Split(humanTime, " ")[0] {
				case "on":
					pieces := strings.Split(humanTime, " ")
					if len(pieces) > 1 {
						if ts, err := dateparse.ParseAny(strings.Join(pieces[1:], " ")); err == nil {
							story.Timestamp = ts
						}
					}
				default:
					if ts := ago.Time(humanTime, now); ts != nil {
						story.Timestamp = *ts
					}
				}

				stories = append(stories, story)
			})
			moreLink = doc.Find(".morelink").Last().AttrOr("href", "")
			if len(moreLink) > 0 && !strings.HasPrefix(moreLink, "https://") {
				moreLink = fmt.Sprintf("%s/%s", baseURL, moreLink)
			}

			if MaxStories != -1 && len(stories) >= MaxStories {
				break
			}
		}

		switch OutputFormat {
		case "json":
			bs, err := json.MarshalIndent(stories, "", "    ")
			if err != nil {
				errorExit(err)
			}
			fmt.Print(string(bs))

		case "yaml":
			bs, err := yaml.Marshal(stories)
			if err != nil {
				errorExit(err)
			}
			fmt.Print(string(bs))

		default:
			errorExit(fmt.Errorf("unrecognized output format %q", OutputFormat))
		}
	},
}

func login(username string, password string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %s", err)
	}

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	form := url.Values{}
	form.Add("acct", username)
	form.Add("pw", password)
	form.Add("goto", "news")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", baseURL), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("login: creating POST request: %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", baseURL+"/")
	req.Header.Add("Origin", baseURL+"/")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Authority", strings.Split(baseURL, "://")[1])
	req.Header.Add("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login: %s", err)
	}
	if resp.StatusCode/100 != 3 || resp.Header.Get("Location") == "" {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("login: expected 3xx reponse status-code but got %v (body=%v)", resp.StatusCode, string(body))
	}
	return client, nil
}

func getLoggedInPage(client *http.Client, page string) (io.ReadCloser, error) {
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

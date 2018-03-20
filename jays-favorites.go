package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/gigawattio/ago"
)

const baseURL = "https://news.ycombinator.com"

var (
	numExpr    = regexp.MustCompile(`^([0-9]+).*$`)
	hnuserExpr = regexp.MustCompile(`^user\?id=`)
)

type Story struct {
	ID          int64
	Title       string
	URL         string
	Points      int64
	Comments    int64
	CommentsURL string
	Submitter   string
	Timestamp   time.Time
}

func main() {
	var (
		stories  = []Story{}
		moreLink = "https://news.ycombinator.com/favorites?id=jaytaylor"
	)

	for len(moreLink) > 0 {
		now := time.Now()

		doc, err := goquery.NewDocument(moreLink)
		if err != nil {
			errorExit(err)
		}

		doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
			var (
				title    = s.Find(".title a")
				comments = s.Next().Find("a").Last()
				story    = Story{
					ID:          int64Or(s.AttrOr("id", "0"), -1),
					Title:       title.Text(),
					URL:         title.AttrOr("href", ""),
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
	}

	bs, err := json.MarshalIndent(stories, "", "    ")
	if err != nil {
		errorExit(err)
	}
	fmt.Print(string(bs))
}

func int64Or(s string, or int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return or
	}
	return i64
}

func errorExit(err interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/gigawattio/ago"
	"github.com/jaytaylor/hn-utils/domain"
)

var (
	numExpr    = regexp.MustCompile(`^([0-9]+).*$`)
	hnuserExpr = regexp.MustCompile(`^user\?id=`)
)

func ExtractStory(s *goquery.Selection) domain.Story {
	var (
		title    = s.Find(".title a.storylink")
		comments = s.Next().Find("a").Last()
		story    = domain.Story{
			ID:          Int64Or(s.AttrOr("id", "0"), -1),
			Title:       title.Text(),
			URL:         ReconstructHNURL(title.AttrOr("href", "")),
			Points:      Int64Or(numExpr.ReplaceAllString(s.Next().Find(".score").Text(), "$1"), -1),
			Comments:    Int64Or(numExpr.ReplaceAllString(comments.Text(), "$1"), -1),
			CommentsURL: comments.AttrOr("href", ""),
			Submitter:   hnuserExpr.ReplaceAllString(s.Next().Find(".hnuser").AttrOr("href", ""), ""),
		}
	)
	if len(story.CommentsURL) > 0 && !strings.HasPrefix(story.CommentsURL, "https://") {
		story.CommentsURL = fmt.Sprintf("%s/%s", BaseURL, story.CommentsURL)
	}

	humanTime := s.Next().Find(".age").Text()
	// For items favorited in the early days of the feature, HN spits out
	// "on <date>" instead of a human-style delta.
	switch strings.Split(humanTime, " ")[0] {
	case "on":
		humanTime = strings.Replace(humanTime, "Sept", "Sep", -1)
		pieces := strings.Split(humanTime, " ")
		if len(pieces) > 1 {
			if ts, err := dateparse.ParseAny(strings.Join(pieces[1:], " ")); err == nil {
				story.Timestamp = ts
			}
		}
	default:
		now := time.Now()
		if ts := ago.Time(humanTime, now); ts != nil {
			story.Timestamp = *ts
		}
	}

	return story
}

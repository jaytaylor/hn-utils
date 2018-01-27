package common

import (
	"fmt"
	"strings"
	"time"

	"gigawatt.io/ago"
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"jaytaylor.com/html2text"

	"jaytaylor.com/hn-utils/domain"
)

// TODO: Add parsing for story root comments (e.g. "Ask HN").

// ExtractDiscussion consumes an HN "/item?id=xxx" page DOM (or subset thereof)
// and parses out all the conversation threads, returning a tree-like
// representation of the entire discussion.
//
// If you have a *goquery.Document, simply pass it in via: doc.Selection.
func ExtractDiscussion(doc *goquery.Selection) domain.Threads {
	var (
		comments = domain.Threads{}
		parents  = map[int]*domain.Comment{} // Track each parent comment at depth X.
	)

	// markParent clears out any parent comments at the same or deeper depth as
	// the passed comment before marking the new parent at depth X.
	markParent := func(c *domain.Comment) {
		if c.Width == 0 {
			parents = map[int]*domain.Comment{}
		} else {
			for x, _ := range parents {
				if x >= c.Width {
					delete(parents, x)
				}
			}
		}
		parents[c.Width] = c
	}

	// findParent locates a comments parent based on it's nesting depth.
	// Nil is returned if no parent is found.
	findParent := func(c *domain.Comment) *domain.Comment {
		if c.Width == 0 {
			return nil
		}
		if p, ok := parents[(c.Depth()-1)*domain.CommentNestingWidthIncrement]; ok {
			return p
		}
		return nil
	}

	doc.Find(".athing.comtr").Each(func(_ int, s *goquery.Selection) {
		c := extractComment(s)
		if c == nil {
			return
		}
		if c.Width == 0 {
			comments = append(comments, c)
		} else if p := findParent(c); p != nil {
			p.Children = append(p.Children, c)
		} else {
			panic(fmt.Errorf("no parent found for comment=%+v", c))
		}
		markParent(c)
	})

	return comments
}

// extractComment returns nil if no comment was found, or the ID or width parse
// fails due to an invalid value.
func extractComment(s *goquery.Selection) *domain.Comment {
	// Remove the "reply" link node from the comment.
	s.Find(".commtext").Children().Each(func(_ int, s *goquery.Selection) {
		if s.HasClass("reply") {
			s.Remove()
		}
	})
	content, _ := html2text.FromHTMLNode(s.Find(".commtext").Get(0))

	c := &domain.Comment{
		ID:        Int64Or(s.AttrOr("id", "0"), -1),
		Author:    s.Find("a.hnuser").First().Text(),
		Timestamp: parseAge(s.Find(".age").Text()),
		Content:   content,
		N:         int(Int64Or(s.Find(".togg").AttrOr("n", "0"), -1)),
		Width:     int(Int64Or(s.Find("img").First().AttrOr("width", "0"), -1)),
	}

	if c != nil && (c.ID == -1 || c.N == -1 || c.Width == -1) {
		c = nil
	}
	return c
}

// parseAge converts any HN time string into a Go time.Time struct.
func parseAge(humanTime string) time.Time {
	// For items favorited in the early days of the feature, HN spits out
	// "on <date>" instead of a human-style delta.
	switch strings.Split(humanTime, " ")[0] {
	case "on":
		humanTime = strings.Replace(humanTime, "Sept", "Sep", -1)
		pieces := strings.Split(humanTime, " ")
		if len(pieces) > 1 {
			if ts, err := dateparse.ParseAny(strings.Join(pieces[1:], " ")); err == nil {
				return ts
			}
		}
	default:
		now := time.Now()
		if ts := ago.Time(humanTime, now); ts != nil {
			return *ts
		}
	}

	return time.Time{}
}

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jaytaylor/hn-utils/common"
	"github.com/jaytaylor/hn-utils/common/storiesflags"
	"github.com/jaytaylor/hn-utils/domain"

	"gigawatt.io/ago"
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	numExpr    = regexp.MustCompile(`^([0-9]+).*$`)
	hnuserExpr = regexp.MustCompile(`^user\?id=`)
)

func init() {
	storiesflags.Init(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		common.ErrorExit(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "favorites",
	Short: "Downloads HN user favorites",
	Long:  "Retrieves HN user favorites as an array of structured Story objects",
	PreRun: func(_ *cobra.Command, _ []string) {
		common.InitLogging(storiesflags.Quiet, storiesflags.Verbose)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			stories  = []domain.Story{}
			moreLink = fmt.Sprintf("%v/favorites?id=%v", common.BaseURL, storiesflags.User)
		)

		for len(moreLink) > 0 {
			now := time.Now()

			log.WithField("more-link", moreLink).Debug("Fetching")
			doc, err := goquery.NewDocument(moreLink)
			if err != nil {
				common.ErrorExit(err)
			}

			doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
				var (
					title    = s.Find(".title a.storylink")
					comments = s.Next().Find("a").Last()
					story    = domain.Story{
						ID:          common.Int64Or(s.AttrOr("id", "0"), -1),
						Title:       title.Text(),
						URL:         common.ReconstructHNURL(title.AttrOr("href", "")),
						Points:      common.Int64Or(numExpr.ReplaceAllString(s.Next().Find(".score").Text(), "$1"), -1),
						Comments:    common.Int64Or(numExpr.ReplaceAllString(comments.Text(), "$1"), -1),
						CommentsURL: comments.AttrOr("href", ""),
						Submitter:   hnuserExpr.ReplaceAllString(s.Next().Find(".hnuser").AttrOr("href", ""), ""),
					}
				)
				if len(story.CommentsURL) > 0 && !strings.HasPrefix(story.CommentsURL, "https://") {
					story.CommentsURL = fmt.Sprintf("%s/%s", common.BaseURL, story.CommentsURL)
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
				moreLink = fmt.Sprintf("%s/%s", common.BaseURL, moreLink)
			}

			if storiesflags.MaxStories != -1 && len(stories) >= storiesflags.MaxStories {
				break
			}
		}

		switch storiesflags.OutputFormat {
		case "json":
			bs, err := json.MarshalIndent(stories, "", "    ")
			if err != nil {
				common.ErrorExit(err)
			}
			fmt.Print(string(bs))

		case "yaml":
			bs, err := yaml.Marshal(stories)
			if err != nil {
				common.ErrorExit(err)
			}
			fmt.Print(string(bs))

		default:
			common.ErrorExit(fmt.Errorf("unrecognized output format %q", storiesflags.OutputFormat))
		}
	},
}

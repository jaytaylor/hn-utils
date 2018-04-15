package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

const baseURL = "https://news.ycombinator.com"

var (
	numExpr    = regexp.MustCompile(`^([0-9]+).*$`)
	hnuserExpr = regexp.MustCompile(`^user\?id=`)
)

var (
	User         string
	OutputFormat string
	MaxStories   int
	Quiet        bool
	Verbose      bool
)

func init() {
	initFavoritesFlags()

	initUpvotesFlags()

	rootCmd.AddCommand(upvotesCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorExit(err)
	}
}

func initFavoritesFlags() {
	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", "jaytaylor", "HN username to find favorites for")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "json", `Output format, one of "json", "yaml"`)
	rootCmd.PersistentFlags().IntVarP(&MaxStories, "max-stories", "m", -1, "Maximum number of stories to collect")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")
}

var rootCmd = &cobra.Command{
	Use:   "favorites",
	Short: "HN user favorites",
	Long:  "Retrieves favorites for a given HN user",
	PreRun: func(_ *cobra.Command, _ []string) {
		initLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			stories  = []domain.Story{}
			moreLink = fmt.Sprintf("%v/favorites?id=%v", baseURL, User)
		)

		for len(moreLink) > 0 {
			now := time.Now()

			log.WithField("more-link", moreLink).Debug("Fetching")
			doc, err := goquery.NewDocument(moreLink)
			if err != nil {
				errorExit(err)
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

func int64Or(s string, or int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return or
	}
	return i64
}

func reconstructHNURL(u string) string {
	if strings.HasPrefix(u, "item?id=") {
		return fmt.Sprintf("%v/%v", baseURL, u)
	}
	return u
}

func errorExit(err interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

func initLogging() {
	level := log.InfoLevel
	if Verbose {
		level = log.DebugLevel
	}
	if Quiet {
		level = log.ErrorLevel
	}
	log.SetLevel(level)
}

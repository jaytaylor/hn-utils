package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jaytaylor/hn-utils/common"
	"github.com/jaytaylor/hn-utils/domain"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	favoritesCmd.MarkFlagRequired("user")
}

var favoritesCmd = &cobra.Command{
	Use:   "favorites",
	Short: "Downloads HN user favorite stories",
	Long:  "Retrieves HN user favorite stories as an array of structured Story objects",
	Args:  cobra.ExactArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		if User == "" || Password == "" {
			log.Warnf("-u/--user and/or -p/--password flag is absent; there is an increased change this client will be blacklisted")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			user            = args[0]
			moreLink        = fmt.Sprintf("%v/favorites?id=%v", common.BaseURL, user)
			stories         = domain.Stories{}
			existingStories domain.Stories
			existingID      int64 = -1 // Used for picking up where an existing collection ends.
			client          *http.Client
			err             error
		)

		if User == "" || Password == "" {
			client = common.NoAuthClient()
		} else {
			if client, err = common.Login(User, Password); err != nil {
				log.Fatal(err)
			}
			log.Debug("Logged in successfully")
		}

		if ReadExisting != "" {
			var err error
			if existingStories, err = common.LoadStories(ReadExisting); err != nil {
				log.Fatal(err)
			} else if len(existingStories) > 0 {
				existingID = existingStories[0].ID
			}
		}

		for len(moreLink) > 0 {
			log.WithField("more-link", moreLink).Debug("Fetching")
			rc, err := common.CheckedGet(client, moreLink)
			if err != nil {
				log.Fatal(err)
			}

			doc, err := goquery.NewDocumentFromReader(rc)
			if err != nil {
				log.Fatal(err)
			}
			if err := rc.Close(); err != nil {
				log.Fatalf("closing response body from %v: %s", moreLink, err)
			}

			var caughtUp bool

			doc.Find(".athing").EachWithBreak(func(i int, s *goquery.Selection) bool {
				if caughtUp {
					return false
				}

				story := common.ExtractStory(s)

				if story.ID == existingID {
					log.WithField("story-id", story.ID).Debug("Caught up to newest story in pre-existing data")
					caughtUp = true
					stories = append(stories, existingStories...)
				}

				stories = append(stories, story)

				return true
			})

			if caughtUp {
				break
			}

			moreLink = doc.Find(".morelink").Last().AttrOr("href", "")
			if len(moreLink) > 0 && !strings.HasPrefix(moreLink, "https://") {
				moreLink = fmt.Sprintf("%s/%s", common.BaseURL, moreLink)
			}

			if MaxItems != -1 && len(stories) >= MaxItems {
				break
			}
		}

		switch OutputFormat {
		case "json":
			bs, err := json.MarshalIndent(stories, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(bs))

		case "yaml":
			bs, err := yaml.Marshal(stories)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(bs))

		default:
			log.Fatalf("unrecognized output format %q", OutputFormat)
		}
	},
}

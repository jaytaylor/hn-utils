package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/jaytaylor/hn-utils/common"
	"github.com/jaytaylor/hn-utils/domain"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	ID           string
	MaxStories   int
	OutputFormat string
	Password     string
	ReadExisting string
	Section      string
	Quiet        bool
	User         string
	Verbose      bool

	// TODO: Add "comments", "story", but will require updates to support
	//       threaded structure.
	Sections = map[string]string{
		"ask": "/ask",
		//"comments":    "/threads?id=%v",
		"favorites": "/favorites?id=%v",
		"frontpage": "/",
		//"story":        "/item?id=%v",
		"new":         "/newest",
		"show":        "/show",
		"submissions": "/submitted?id=%v",
		"upvotes":     "/upvoted?id=%v",
	}
)

func init() {
	areas := []string{}
	for s, _ := range Sections {
		areas = append(areas, s)
	}

	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", "jaytaylor", "HN username to login as")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "HN login password")
	rootCmd.PersistentFlags().StringVarP(&ID, "id", "i", "", "Relevant user or story identifier")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "json", `Output format, one of "json", "yaml"`)
	rootCmd.PersistentFlags().IntVarP(&MaxStories, "max-stories", "m", -1, "Maximum number of stories to collect")
	rootCmd.PersistentFlags().StringVarP(&ReadExisting, "existing", "e", "", `Load an existing array of stories from named JSON database file, then front-load new content (set to "-" to read from STDIN)`)
	rootCmd.PersistentFlags().StringVarP(&Section, "section", "s", "frontpage", fmt.Sprintf("Site area to get paged results for.  Available selections: %v", strings.Join(areas, ", ")))
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hn-slurp",
	Short: "Download the specified section from HN and transform it into structured JSON",
	Long:  "Retrieves objects as an array of structured Story object for a given HN user/password combination.  The 'user upvotes' section has a hard requirement for user/password login.",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		common.InitLogging(Quiet, Verbose)

		// Validate section.
		var found bool
		for s, _ := range Sections {
			if strings.ToLower(Section) == s {
				found = true
				break
			}
		}
		if !found {
			return errors.New("Missing required flag: -s/--section must not be empty, see --help for a lis of valid secitions")
		}

		// Require password supplied when collecting user upvotes.
		if Section == "upvotes" && Password == "" {
			return errors.New("Missing required flag: -p/--password must not be empty")
		}

		// Validate ID.
		if strings.Contains(Sections[Section], "%v") {
			if ID == "" {
				return fmt.Errorf("Missing required flag: -i/--id must not be empty for section=%v", Section)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			moreLink        = fmt.Sprintf("%v%v", common.BaseURL, Sections[Section])
			stories         = domain.Stories{}
			existingStories domain.Stories
			existingID      int64 = -1 // Used for picking up where an existing collection ends.
		)

		if strings.Contains(moreLink, "%v") {
			// Fill in ID param.
			moreLink = fmt.Sprintf(moreLink, ID)
		}

		if ReadExisting != "" {
			var err error
			if existingStories, err = common.LoadStories(ReadExisting); err != nil {
				return err
			} else if len(existingStories) > 0 {
				existingID = existingStories[0].ID
			}
		}

		client, err := getClient()
		if err != nil {
			return err
		}
		log.Debug("Logged in successfully")

		for len(moreLink) > 0 {
			log.WithField("more-link", moreLink).Debug("Fetching")

			rc, err := common.CheckedGet(client, moreLink)
			if err != nil {
				return err
			}

			doc, err := goquery.NewDocumentFromReader(rc)
			if err != nil {
				return err
			}
			if err := rc.Close(); err != nil {
				return fmt.Errorf("closing response body from %v: %s", moreLink, err)
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

			if MaxStories != -1 && len(stories) >= MaxStories {
				break
			}
		}

		switch OutputFormat {
		case "json":
			bs, err := json.MarshalIndent(stories, "", "    ")
			if err != nil {
				return err
			}
			fmt.Print(string(bs))

		case "yaml":
			bs, err := yaml.Marshal(stories)
			if err != nil {
				return err
			}
			fmt.Print(string(bs))

		default:
			return fmt.Errorf("unrecognized output format %q", OutputFormat)
		}

		return nil
	},
}

func getClient() (*http.Client, error) {
	if Password == "" {
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
		return client, nil
	}

	client, err := common.Login(User, Password)
	return client, err
}

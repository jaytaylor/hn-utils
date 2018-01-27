package main

import (
	"github.com/jaytaylor/hn-utils/common"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const defaultUser = "jaytaylor"

var (
	Quiet        bool
	Verbose      bool
	User         string
	Password     string
	OutputFormat string
	MaxItems     int
	ReadExisting string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")

	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", defaultUser, "HN username to authenticate with")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "HN account password")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "json", `Output format, one of: "json", "yaml"`)
	rootCmd.PersistentFlags().IntVarP(&MaxItems, "max", "m", -1, "Maximum number of items to collect (when applicable)")
	rootCmd.PersistentFlags().StringVarP(&ReadExisting, "existing", "e", "", `Load an existing array of items from named JSON database file and front-load new content (set to "-" to read from STDIN)`)

	rootCmd.AddCommand(
		favoritesCmd,
		upvotedCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hn",
	Short: "HN data retrieval tools",
	Long:  "Tools for retrieving data from HackerNews (news.ycombinator.com) via scraping",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		common.InitLogging(Quiet, Verbose)
	},
}

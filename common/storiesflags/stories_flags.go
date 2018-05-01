package storiesflags

import (
	"github.com/spf13/cobra"
)

var (
	User         string
	OutputFormat string
	MaxStories   int
	EditInPlace  string
	Quiet        bool
	Verbose      bool
)

// Init adds common stories-related command-line flags to a *cobra.Command.
func Init(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&User, "user", "u", "jaytaylor", "HN username to find favorites for")
	cmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "json", `Output format, one of "json", "yaml"`)
	cmd.PersistentFlags().IntVarP(&MaxStories, "max-stories", "m", -1, "Maximum number of stories to collect")
	cmd.PersistentFlags().StringVarP(&EditInPlace, "edit-in-place", "i", "", `Update an existing JSON database file with new content (or set to "-" to use STDIN / STDOUT)`)
	cmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")
}

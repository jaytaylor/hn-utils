package storiesflags

import (
	"github.com/spf13/cobra"
)

var (
	User         string
	OutputFormat string
	MaxStories   int
	ReadExisting string
	Quiet        bool
	Verbose      bool
)

// Init adds common stories-related command-line flags to a *cobra.Command.
func Init(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&User, "user", "u", "jaytaylor", "HN username to find favorites for")
	cmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "json", `Output format, one of "json", "yaml"`)
	cmd.PersistentFlags().IntVarP(&MaxStories, "max-stories", "m", -1, "Maximum number of stories to collect")
	cmd.PersistentFlags().StringVarP(&ReadExisting, "existing", "e", "", `Load an existing array of Stories from named JSON database file and front-load new content (set to "-" to read from STDIN)`)
	cmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "Activate quiet log output")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Activate verbose log output")
}

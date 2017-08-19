package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
)

// fullAPIKey flag for troubleshoot command.
var fullAPIKey bool

// troubleshootCmd does a diagnostic self-check.
var troubleshootCmd = &cobra.Command{
	Use:     "troubleshoot",
	Aliases: []string{"t"},
	Short:   "Troubleshoot does a diagnostic self-check.",
	Long: `Provides output to help with troubleshooting.

If you're running into trouble, copy and paste the output from the troubleshoot
command into a GitHub issue so we can help figure out what's going on.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.HTTPClient = &http.Client{Timeout: 20 * time.Second}
		c := cli.New(Version)

		cfg, err := config.NewUserConfig()
		BailOnError(err)

		status := cli.NewStatus(c, *cfg)
		status.Censor = !fullAPIKey
		s, err := status.Check()
		BailOnError(err)
		fmt.Printf("%s", s)
	},
}

func init() {
	RootCmd.AddCommand(troubleshootCmd)
	troubleshootCmd.Flags().BoolVarP(&fullAPIKey, "full-api-key", "f", false, "display the user's full API key, censored by default")
}

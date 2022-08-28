package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
)

// versionString returns the version prefixed by 'v'
// or an empty string if no version has been populated by goreleaser.
// In this case, the --version flag will not be added by cobra.
func versionString() string {
	if len(version) == 0 {
		return "0.0.0"
	}
	return version
}

func NewCmdVersion() *cobra.Command {
	// certCmd represents the cert command
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version for cm",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(versionString())
		},
	}

	return cmd
}

package hx

import (
	"os"

	"github.com/bryant-rh/cm/pkg/buildx"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdHxBuildXCISetup() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "ci-setup",
		Short: "buildx setup in ci",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infoln("hx ci-setup")
		},
		Run: func(cmd *cobra.Command, args []string) {
			buildkitName := os.Getenv("BUILDKIT_NAME")
			if buildkitName == "" {
				buildkitName = "buildkit"
			}

			if len(args) > 0 {
				buildkitName = args[0]
			}
			buildx.RunDockerBuildxSetup(buildkitName)
		},
	}

	return cmd
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// buildKernel sets up the environment, creates the output directory and builds
// the kernel.
func buildKernel(cmd *cobra.Command, args []string) {
	var err error

	profile, err := getBuildConf(cmd)
	errFatal(err)

	err = profile.Setup()
	errFatal(err)

	err = os.MkdirAll(profile.BuildDir, 0755)
	errFatal(err)

	err = profile.Config()
	errFatal(err)
}

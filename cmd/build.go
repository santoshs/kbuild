package cmd

import (
	"fmt"
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

	// If the user provides a build command explicitly only execute that
	if len(args) > 0 {
		errFatal(profile.Build(args))
		return
	}

	if skipconfig, err := cmd.Flags().GetBool("skip-config"); err != nil {
		errFatal(err)
	} else if !skipconfig {
		errFatal(profile.Config())
	}

	errFatal(profile.Build(args))
	for _, m := range profile.Modules {
		if err := profile.Build([]string{
			fmt.Sprintf("M=%s", m), "modules"}); err != nil {
			errFatal(err)
		}
	}
}

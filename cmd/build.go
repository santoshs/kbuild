package cmd

import (
	"github.com/spf13/cobra"

	"github.com/santoshs/kbuild/pkg/kbuild"
)

// buildKernel sets up the environment, creates the output directory and builds
// the kernel.
func buildKernel(cmd *cobra.Command, args []string) {
	var err error

	profile, err := getBuildConf(cmd)
	errFatal(err)

	kb, err := kbuild.NewKbuild(profile.SrcPath, profile.BuildPath)

	kb.SetArch(profile.Arch)
	kb.NumParallelJobs = profile.NumJobs
	kb.BuildDir = profile.BuildDir
	kb.Pull = profile.Pull

	err = kb.Build(args)
	errFatal(err)
}

package cmd

import (
	"github.com/spf13/cobra"
)

// kernelBuild ...
func buildKernel(cmd *cobra.Command, args []string) {
	var err error

	kb, err := getkbuild(cmd)
	errFatal(err)

	if getArg(cmd, "dry-run", false).(bool) {
		return
	}

	err = kb.Build(args)
	errFatal(err)
}

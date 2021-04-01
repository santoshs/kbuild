package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// kernelBuild ...
func buildKernel(cmd *cobra.Command, args []string) {
	var err error

	kb, err := getkbuild(cmd)
	if err != nil {
		log.Fatal(err)
	}

	dry_run, err := envArgBool(cmd, "dry-run")
	errFatal(err)
	if dry_run {
		return
	}

	err = kb.Build(args)
	errFatal(err)
}

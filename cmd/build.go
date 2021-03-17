package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// kernelBuild ...
func buildKernel(cmd *cobra.Command, args []string) {
	var err error

	kb, err := getkbuild(cmd, args)
	if err != nil {
		log.Fatal(err)
	}

	if err := kb.Build(); err != nil {
		log.Fatal(err)
	}
}

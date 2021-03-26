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

	if err := kb.Build(args); err != nil {
		log.Fatal(err)
	}
}

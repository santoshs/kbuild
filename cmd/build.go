package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/santoshs/kbuild/pkg/kbuild"
)

// kernelBuild ...
func buildKernel(cmd *cobra.Command, args []string) {
	var source string
	var err error

	if len(args) > 0 {
		source = args[0]
	} else {
		if source, err = os.Getwd(); err != nil {
			log.Fatal(err)
		}
	}

	kb, err := kbuild.NewKbuild(source, "~/.cache/kbuild")
	if err != nil {
		log.Fatal(err)
	}

	arch, err := cmd.Flags().GetString("arch")
	if err == nil {
		kb.Arch = arch
	}

	if err := kb.Build(); err != nil {
		log.Fatal()
	}
}

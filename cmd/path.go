package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show path of different build artefacts",
	Long:  "",
	Run:   showPath,
}

// showPath ...
func showPath(cmd *cobra.Command, args []string) {
	kb, err := getkbuild(cmd)
	if err != nil {
		log.Fatal(err)
	}
	dir, err := kb.GetBuildDir()
	if err != nil {
		log.Fatal(err)
	}

	bz, err := cmd.Flags().GetBool("bzimage")
	if err != nil {
		log.Fatal(err)
	}

	if bz {
		fmt.Println(fmt.Sprintf("%s/arch/%s/boot/bzImage",
			dir, kb.Arch))
		return
	}

	c, err := cmd.Flags().GetBool("config")
	if err != nil {
		log.Fatal(err)
	}
	if c {
		fmt.Println(fmt.Sprintf("%s/.config", dir))
		return
	}

	fmt.Println(dir)
}

func init() {
	pathCmd.Flags().BoolP("bzimage", "z", false, "Show bzimage path")
	pathCmd.Flags().BoolP("config", "c", false, "Show .config path")
}
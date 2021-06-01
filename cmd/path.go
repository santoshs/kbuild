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
	profile, err := getBuildConf(cmd)
	errFatal(err)

	profile.Setup()

	if bz, err := cmd.Flags().GetBool("bzimage"); err != nil {
		log.Fatal(err)
	} else if bz {
		fmt.Println(fmt.Sprintf("%s/arch/%s/boot/bzImage",
			profile.BuildDir, profile.Arch))
		return
	}

	if c, err := cmd.Flags().GetBool("config"); err != nil {
		log.Fatal(err)
	} else if c {
		fmt.Println(fmt.Sprintf("%s/.config", profile.BuildDir))
		return
	}

	if v, err := cmd.Flags().GetBool("vmlinux"); err != nil {
		log.Fatal(err)
	} else if v {
		fmt.Println(fmt.Sprintf("%s/vmlinux", profile.BuildDir))
		return
	}

	fmt.Println(profile.BuildDir)
}

func init() {
	pathCmd.Flags().BoolP("bzimage", "z", false, "Show bzimage path")
	pathCmd.Flags().BoolP("config", "c", false, "Show .config path")
	pathCmd.Flags().BoolP("vmlinux", "l", false, "Show vmlinux path")
}

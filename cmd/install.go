package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/santoshs/kbuild/pkg/kbuild"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install kernel",
	Long:  "",
	Run:   installKernel,
}

// installKernel ...
func installKernel(cmd *cobra.Command, args []string) {
	var err error
	var install_kernel = true
	install_path := ""

	profile, err := getBuildConf(cmd)
	errFatal(err)

	kb, err := kbuild.NewKbuild(profile.SrcPath, profile.BuildDir)
	errFatal(err)

	image, err := cmd.Flags().GetString("image")
	errFatal(err)

	if image != "" {
		if s, err := os.Stat(image); err == nil {
			if s.IsDir() {
				errFatal(fmt.Errorf("Image cannot be a directory\n"))
			}
		} else {
			errFatal(err)
		}

		// create a temporary mount point
		mp, err := ioutil.TempDir("", "kbuild")
		errFatal(err)

		errFatal(nbdMount(image, mp))
		defer nbdUmount(mp)

		install_path = mp
	}

	if _, err = cmd.Flags().GetBool("modules-only"); err != nil {
		install_kernel = false
	}

	err = kb.Install(install_path, true, install_kernel, args)
	errFatal(err)
}

func init() {
	installCmd.Flags().StringP("image", "q", "",
		"Install in the given qemu image")
	installCmd.Flags().BoolP("modules-only", "M", false,
		"install only modules")
}

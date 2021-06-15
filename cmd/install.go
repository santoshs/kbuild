package cmd

import (
	"fmt"

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

	profile, err := getBuildConf(cmd)
	errFatal(err)
	profile.Setup()

	ipath, err := cmd.Flags().GetString("install-path")
	errFatal(err)
	if ipath != "" {
		profile.Environment["INSTALL_PATH"] = ipath
	}

	mpath, err := cmd.Flags().GetString("install-mod-path")
	errFatal(err)
	if mpath != "" {
		profile.Environment["INSTALL_MOD_PATH"] = mpath
	}

	// If there are options module build directories are provided we will
	// install those too.
	for _, m := range profile.Modules {
		if err := runCmd("sudo", []string{"-E", "--", "make",
			fmt.Sprintf("M=%s", m), "modules_install"},
			profile.getenv()); err != nil {
			errFatal(err)
		}
	}

	// This should be last, because this will do the depmod
	errFatal(runCmd("sudo", []string{"-E", "--", "make", "modules_install"}, profile.getenv()))

	if m, err := cmd.Flags().GetBool("modules-only"); err != nil {
		errFatal(err)
	} else if m {
		return
	}

	errFatal(runCmd("sudo", []string{"-E", "--", "make", "install"}, profile.getenv()))
}

func init() {
	installCmd.Flags().StringP("install-path", "i", "",
		"Install in the given path")
	installCmd.Flags().StringP("install-mod-path", "m", "",
		"Install modules in the given path")
	installCmd.Flags().BoolP("modules-only", "M", false,
		"install only modules")
}

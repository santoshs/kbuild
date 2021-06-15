package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbuild",
	Short: "Kernel build helper",
	Long:  "",
	Args:  cobra.ArbitraryArgs,
	Run:   buildKernel,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cwd, err := os.Getwd()
	errFatal(err)

	rootCmd.PersistentFlags().IntP("jobs", "j", runtime.NumCPU()/2,
		"Number of jobs to run")
	rootCmd.PersistentFlags().CountP("verbose", "v",
		"Verbose output, the more the 'v's the more verbose")
	rootCmd.PersistentFlags().StringP("arch", "a", GetHostArch(),
		"Target architecture")
	rootCmd.PersistentFlags().StringP("builddir", "o", "",
		`Name of the build directory. Can also be set using
KBUILD_BUILDDIR environment variable. (default: ~/.cache/kbuild/srcdir.branch.arch)`)
	rootCmd.PersistentFlags().StringP("srcdir", "s", cwd,
		"Path to the source directory, defaults to current directory")
	rootCmd.PersistentFlags().Bool("pull", false,
		"Update the source repository")
	rootCmd.PersistentFlags().Bool("dry-run", false,
		"For debugging; do not do anything")
	rootCmd.PersistentFlags().StringP("profile", "p", "",
		`Use the specified profile from the config file. Individual
config items can be overridden through the CLI arguments
or environment variables`)
	rootCmd.Flags().MarkHidden("dry-run")
	rootCmd.Flags().BoolP("skip-config", "S", false,
		`Do not make config, skip merging configure fragments too
if present in profile`)
	rootCmd.Flags().BoolP("skip-base", "B", false, "Do not make defconfig")

	rootCmd.AddCommand(pathCmd)
	rootCmd.AddCommand(installCmd)
}

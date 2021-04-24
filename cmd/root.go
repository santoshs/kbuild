package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/santoshs/kbuild/pkg/kbuild"
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
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().IntP("jobs", "j", runtime.NumCPU()/2,
		"Number of jobs to run")
	rootCmd.PersistentFlags().CountP("verbose", "v",
		"Verbose output, the more the 'v's the more verbose")
	rootCmd.PersistentFlags().StringP("arch", "a", kbuild.GetHostArch(),
		"Target architecture")
	rootCmd.PersistentFlags().StringP("buildpath", "b", "~/.cache/kbuild",
		"Build path")
	rootCmd.PersistentFlags().StringP("builddir", "o", "",
		`Name of the build directory, cannot be a path. Can also be
set using KBUILD_BUILD_DIR environment variable`)
	rootCmd.PersistentFlags().StringP("srcdir", "s", cwd,
		"Path to the source directory")
	rootCmd.PersistentFlags().Bool("pull", false,
		"Update the source repository")
	rootCmd.PersistentFlags().Bool("dry-run", false,
		"For debugging; do not do anything")
	rootCmd.Flags().MarkHidden("dry-run")

	rootCmd.AddCommand(pathCmd)
	rootCmd.AddCommand(installCmd)
}

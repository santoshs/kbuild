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
	Args:  cobra.MaximumNArgs(1),
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

	rootCmd.Flags().IntP("jobs", "j", runtime.NumCPU()/2,
		"Number of jobs to run")
	rootCmd.PersistentFlags().CountP("verbose", "v",
		"Verbose output, the more the 'v's the more verbose")
	rootCmd.PersistentFlags().StringP("arch", "a", kbuild.GetHostArch(),
		"Target architecture")
	rootCmd.PersistentFlags().StringP("buildpath", "p", "~/.cache/kbuild",
		"Build path")
	rootCmd.PersistentFlags().StringP("builddir", "b", "",
		"Name of the build directory, cannot be a path")
	rootCmd.PersistentFlags().StringP("srcdir", "s", cwd,
		"Path to the source directory")

	rootCmd.AddCommand(pathCmd)
	pathCmd.Flags().BoolP("bzimage", "z", false, "Show bzimage path")
	pathCmd.Flags().BoolP("config", "c", false, "Show .config path")
}

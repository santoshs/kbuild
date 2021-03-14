package cmd

import (
	"runtime"

	"github.com/santoshs/kbuild/pkg/kbuild"
	"github.com/spf13/cobra"
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
	rootCmd.Flags().CountP("verbose", "v",
		"Verbose output, the more the 'v's the more verbose")
	rootCmd.Flags().IntP("jobs", "j", runtime.NumCPU()/2, "Number of jobs to run")
	rootCmd.Flags().StringP("arch", "a", kbuild.GetHostArch(), "Target architecture")
}

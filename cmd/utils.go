package cmd

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/santoshs/kbuild/pkg/kbuild"
)

type BuildConf struct {
	SrcPath      string `yaml:"srcdir"`
	BuildDir     string `yaml:"builddir"`
	BuildPath    string `yaml:"buildpath"`
	WorktreePath string `yaml:"worktreepath"`
	DefArch      string `yaml:"arch"`
	ToolChain    string `yaml:"toolchain"`
	CCPrefix     string `yaml:"ccprefix"`
}

const BUILD_PATH = "~/.cache/kbuild"

func getkbuild(cmd *cobra.Command) (*kbuild.Kbuild, error) {
	var err error

	buildpath, _ := cmd.Flags().GetString("buildpath")
	src, _ := cmd.Flags().GetString("srcdir")

	kb, err := kbuild.NewKbuild(src, buildpath)
	errFatal(err)

	kb.Arch, err = envArgString(cmd, "arch")
	errFatal(err)

	kb.NumParallelJobs, err = envArgInt(cmd, "jobs")
	errFatal(err)

	kb.BuildDir, err = envArgString(cmd, "builddir")
	errFatal(err)

	kb.NoPull, err = envArgBool(cmd, "no-pull")
	errFatal(err)

	return kb, nil
}

func errFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func errLog(err error) {
	if err != nil {
		log.Println(err)
	}
}

func nbdMount(image, mountpoint string) error {
	// os.Command("")
	return nil
}

func nbdUmount(mountpoint string) error {
	return nil
}

// For all the following functions, related to getting arguments from the
// environment or the command line, first check the environment if we have a
// default, because we should override command line defaults if an arg is set in
// the environment

func envArgString(cmd *cobra.Command, arg string) (string, error) {
	var s string

	s = os.Getenv("KBUILD_" + strings.ToUpper(arg))
	if s != "" {
		return s, nil
	}

	return cmd.Flags().GetString(arg)
}

func envArgInt(cmd *cobra.Command, arg string) (int, error) {
	var s string

	// first check the environment if we have a default, because we should
	// override command line defaults if an arg is set in the environment
	s = os.Getenv("KBUILD_" + strings.ToUpper(arg))
	if s != "" {
		return strconv.Atoi(s)
	}

	return cmd.Flags().GetInt(arg)
}

func envArgBool(cmd *cobra.Command, arg string) (bool, error) {
	var s string

	// first check the environment if we have a default, because we should
	// override command line defaults if an arg is set in the environment
	s = os.Getenv("KBUILD_" + strings.ToUpper(arg))
	if s != "" {
		log.Printf("Using %s from environment\n", arg)
		return true, nil
	}

	return cmd.Flags().GetBool(arg)
}

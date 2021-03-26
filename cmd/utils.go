package cmd

import (
	"log"

	"github.com/santoshs/kbuild/pkg/kbuild"
	"github.com/spf13/cobra"
)

type BuildConf struct {
	SrcPath      string `yaml:"srcdir"`
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
	if err != nil {
		log.Fatal(err)
	}

	if arch, err := cmd.Flags().GetString("arch"); err == nil {
		kb.Arch = arch
	}

	if jobs, err := cmd.Flags().GetInt("jobs"); err == nil {
		kb.NumParallelJobs = jobs
	}

	if dir, err := cmd.Flags().GetString("builddir"); err == nil {
		kb.BuildDir = dir
	}

	return kb, nil
}

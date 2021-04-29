package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Profile struct {
	SrcPath      string            `yaml:"srcdir"`
	BuildDir     string            `yaml:"builddir"`
	Arch         string            `yaml:"arch"`
	CC           string            `yaml:"cc"`
	CrossCompile string            `yaml:"cross_compile"`
	Pull         bool              `yaml:"pull"`
	BaseConfig   string            `yaml:"baseconfig"`
	Configs      []string          `yaml:"configs"`
	Environment  map[string]string `yaml:"env"`
	NumJobs      int               `yaml:"jobs"`

	name string
	repo *git.Repository
	wt   *git.Worktree
}

type KbuildConfig struct {
	Common   *Profile            `yaml:"Common"`
	Profiles map[string]*Profile `yaml:"Profiles"`
}

// setup does some sanity checks and creates the build directory
func (p *Profile) Setup() error {
	var err error

	opts := git.PlainOpenOptions{
		// required to get the Branch() and HEAD() from linked worktrees
		EnableDotGitCommonDir: true,
	}

	p.repo, err = git.PlainOpenWithOptions(p.SrcPath, &opts)
	if err != nil {
		errLog(err)
	} else {
		p.wt, err = p.repo.Worktree()
		if err != nil {
			errLog(err)
		}
	}

	if p.BuildDir == "" {
		p.BuildDir, err = expandHome("~/.cache/kbuild")
		if err != nil {
			return err
		}
		p.mkBuildPath()
	}

	if p.Environment == nil {
		p.Environment = make(map[string]string)
	}
	// TODO make sure that the user specified environment variables are not
	// overridden
	p.Environment["ARCH"] = p.Arch
	p.Environment["KBUILD_OUTPUT"] = p.BuildDir

	return nil
}

// If the user has not given a build directory name, we will create a
// name from src dir name, the branch, architecture and the profile For
// example if the source directory is "linux", branch master, build for
// powerpc64 with profile name profile, the build directory will be
// linux.master.powerpc.profile
func (p *Profile) mkBuildPath() {
	name := ""
	wd := path.Base(p.SrcPath)
	var dir string

	if p.repo != nil {
		ref, err := p.repo.Head()
		if err == nil {
			name = path.Base(ref.Name().String())
		}
	}

	if name != "" {
		dir = fmt.Sprintf("%s.%s.%s.%s", wd, name, p.Arch, p.name)
	} else {
		dir = fmt.Sprintf("%s.%s.%s", wd, p.Arch, p.name)
	}

	p.BuildDir = fmt.Sprintf("%s/%s", p.BuildDir, dir)
}

func (p *Profile) getenv() []string {
	var env []string

	env = os.Environ()
	for v := range p.Environment {
		env = append(env, fmt.Sprintf("%s=%s", v, p.Environment[v]))
	}

	return env
}

func (p *Profile) Config() error {
	err := p.mkConfig()
	if err != nil {
		return err
	}

	err = p.mergeConfig()
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) mkConfig() error {
	if p.BaseConfig == "" {
		p.BaseConfig = "defconfig"
	}

	args := []string{fmt.Sprintf("O=%s", p.BuildDir), p.BaseConfig}

	return runCmd("make", args, p.getenv())
}

func (p *Profile) mergeConfig() error {
	configs := strings.Join(p.Configs, " ")
	if configs == "" {
		return nil
	}

	mergecmd := fmt.Sprintf("%s/scripts/kconfig/merge_config.sh", p.SrcPath)

	args := []string{"-m", configs, fmt.Sprintf("-O=%s", p.BuildDir)}

	return runCmd(mergecmd, args, p.getenv())
}

package cmd

import (
	"fmt"
	"os"
	"path"

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
	Modules      []string          `yaml:"module_paths"`
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

	p.Setenv("ARCH", p.Arch, false)
	p.Setenv("KBUILD_OUTPUT", p.BuildDir, false)

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

func (p *Profile) Config(skip_base bool) error {
	if !skip_base {
		if err := p.mkConfig(); err != nil {
			return err
		}
	}

	if err := p.mergeConfig(); err != nil {
		return err
	}

	return nil
}

func (p *Profile) getArgs() []string {
	// Short options with arguments after a space will be treated as a
	// separate argument, so use long options with arguments after '='.
	return []string{
		fmt.Sprintf("--jobs=%d", p.NumJobs),
	}
}

func (p *Profile) mkConfig() error {
	if p.BaseConfig == "" {
		p.BaseConfig = "defconfig"
	}

	args := []string{p.BaseConfig}
	args = append(args, p.getArgs()...)

	return runCmd("make", args, p.getenv())
}

func (p *Profile) mergeConfig() error {
	var configs []string

	if len(p.Configs) == 0 {
		return nil
	}

	for _, c := range p.Configs {
		cp, err := expandHome(c)
		if err != nil {
			errLog(err)
			continue
		}
		configs = append(configs, cp)
	}

	mergecmd := fmt.Sprintf("%s/scripts/kconfig/merge_config.sh", p.SrcPath)

	args := []string{"-m", fmt.Sprintf("%s/.config", p.BuildDir)}
	args = append(args, configs...)
	env := p.getenv()
	env = append(env, fmt.Sprintf("KCONFIG_CONFIG=%s/.config", p.BuildDir))

	if err := runCmd(mergecmd, args, env); err != nil {
		return err
	}

	// Don't want the build to prompt user for default values again
	return runCmd("make", []string{"olddefconfig"}, env)
}

func (p *Profile) Build(build_args []string) error {
	args := []string{"--append", "--", "make"}

	args = append(args, p.getArgs()...)
	args = append(args, build_args...)

	if err := runCmd("bear", args, p.getenv()); err != nil {
		return err
	}

	return nil
}

// Setenv sets the build environment variables with an option to not overwrite
// existing environment variables.
func (p *Profile) Setenv(key, val string, overwrite bool) {
	if _, ok := p.Environment[key]; ok && !overwrite {
		return
	}

	// we will also have to check the /real/ environment that is set before
	// calling this program, and make sure we don't overwrite what the user
	// has set.
	if os.Getenv(key) != "" && !overwrite {
		return
	}

	p.Environment[key] = val
}

package cmd

import (
	"github.com/go-git/go-git/v5"
)

type Profile struct {
	SrcPath      string            `yaml:"srcdir"`
	BuildDir     string            `yaml:"builddir"`
	BuildPath    string            `yaml:"buildpath"`
	Arch         string            `yaml:"arch"`
	CC           string            `yaml:"cc"`
	CrossCompile string            `yaml:"cross_compile"`
	Pull         bool              `yaml:"pull"`
	BaseConfig   string            `yaml:"baseconfig"`
	Configs      []string          `yaml:"configs"`
	Environment  map[string]string `yaml:"env"`
	NumJobs      int               `yaml:"jobs"`

	repo *git.Repository
	wt   *git.Worktree
}

type KbuildConfig struct {
	Common   *Profile            `yaml:"Common"`
	Profiles map[string]*Profile `yaml:"Profiles"`
}

// setup does some sanity checks and creates the build directory
func (p *Profile) setup() error {
	return nil
}

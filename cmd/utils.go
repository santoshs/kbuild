package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/santoshs/kbuild/pkg/kbuild"
)

type BuildConf struct {
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
}

type KbuildConfig struct {
	Common   *BuildConf            `yaml:"Common"`
	Profiles map[string]*BuildConf `yaml:"Profiles"`
}

const BUILD_PATH = "~/.cache/kbuild"

func loadConf(confFile string) (*KbuildConfig, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	if confFile == "" {
		confFile = filepath.Join(usr.HomeDir, ".config/kbuild")
	}

	f, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	kconf := KbuildConfig{}
	err = yaml.Unmarshal(f, &kconf)

	return &kconf, nil
}

func getkbuild(cmd *cobra.Command) (*kbuild.Kbuild, error) {
	var err error

	_, err = loadConf("")
	errFatal(err)

	buildpath, _ := cmd.Flags().GetString("buildpath")
	src, _ := cmd.Flags().GetString("srcdir")

	kb, err := kbuild.NewKbuild(src, buildpath)
	errFatal(err)

	arch, _ := envArgString(cmd, "arch")
	kb.SetArch(arch)

	kb.NumParallelJobs, err = envArgInt(cmd, "jobs")
	errFatal(err)
	if kb.NumParallelJobs < 1 {
		kb.NumParallelJobs = 1
	}

	kb.BuildDir, err = envArgString(cmd, "builddir")
	errFatal(err)

	kb.Pull, err = envArgBool(cmd, "pull")
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

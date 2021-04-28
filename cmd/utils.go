package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	NumJobs      int               `yaml:"jobs"`

	repo *git.Repository
	wt   *git.Worktree
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

func getBuildConf(cmd *cobra.Command) (*BuildConf, error) {
	var err error
	var profile *BuildConf

	kconf, err := loadConf("")
	errFatal(err)

	pname, _ := getArg(cmd, "profile", "default").(string)
	profile, ok := kconf.Profiles[pname]
	if pname == "default" && !ok {
		kconf.Profiles[pname] = &BuildConf{}
		profile = kconf.Profiles[pname]
	}

	if !ok {
		fmt.Print("Available Profiles: [")
		for k := range kconf.Profiles {
			fmt.Print(k, ",")
		}
		fmt.Println("]")
		errFatal(fmt.Errorf("Profile %s not found", pname))
	}

	profile.BuildPath = getArg(cmd, "buildpath", profile.BuildPath).(string)

	profile.SrcPath = getArg(cmd, "srcdir", profile.SrcPath).(string)
	profile.Arch = getArg(cmd, "arch", profile.Arch).(string)

	profile.NumJobs = getArg(cmd, "jobs", profile.NumJobs).(int)
	if profile.NumJobs < 1 {
		profile.NumJobs = 1
	}

	profile.BuildDir = getArg(cmd, "builddir", profile.BuildDir).(string)

	profile.Pull = getArg(cmd, "pull", profile.Pull).(bool)

	return profile, nil
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
// environment or the command line or the config file. The
//
// 1. Command line argument will override
// 2. Environment variables, which will override
// 3. Setting in the config file, if any.

func getArg(cmd *cobra.Command, arg string, defval interface{}) interface{} {
	var val interface{}
	var err error
	hasEnv := false

	env := os.Getenv("KBUILD_" + strings.ToUpper(arg))
	if env != "" {
		hasEnv = true
	}

	switch defval.(type) {
	case string:
		val, err = cmd.Flags().GetString(arg)
		if defval == "" {
			defval = val
		}
		if hasEnv {
			defval = env
		}
	case int:
		val, err = cmd.Flags().GetInt(arg)
		if defval == 0 {
			defval = val
		}
		if hasEnv {
			defval, err = strconv.Atoi(env)
		}

	case bool:
		val, err = cmd.Flags().GetBool(arg)
		if hasEnv {
			defval = true
		}
	}

	if err != nil || cmd.Flags().Changed(arg) {
		return val
	}

	return defval
}

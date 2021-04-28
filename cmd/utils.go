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
	"syscall"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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

func getBuildConf(cmd *cobra.Command) (*Profile, error) {
	var err error
	var profile *Profile

	kconf, err := loadConf("")
	errFatal(err)

	pname, _ := getArg(cmd, "profile", "default").(string)
	profile, ok := kconf.Profiles[pname]
	if pname == "default" && !ok {
		kconf.Profiles[pname] = &Profile{}
		profile = kconf.Profiles[pname]
		ok = true
	} else if !ok {
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

func GetHostArch() string {
	var arch []byte
	utsname := syscall.Utsname{}
	syscall.Uname(&utsname)

	for _, v := range utsname.Machine {
		if v == 0 {
			break
		}
		arch = append(arch, byte(v))
	}

	return string(arch)
}

func expandHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

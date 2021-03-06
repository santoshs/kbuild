package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

	profile.name = pname

	profile.SrcPath = getArg(cmd, "srcdir", profile.SrcPath).(string)
	profile.SrcPath, err = expandHome(profile.SrcPath)
	errFatal(err)

	profile.Arch = getArg(cmd, "arch", profile.Arch).(string)

	profile.NumJobs = getArg(cmd, "jobs", profile.NumJobs).(int)
	if profile.NumJobs < 1 {
		profile.NumJobs = 1
	}

	profile.BuildDir = getArg(cmd, "builddir", profile.BuildDir).(string)
	if profile.BuildDir != "" {
		profile.BuildDir, err = expandHome(profile.BuildDir)
		errFatal(err)
	}

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

func pipetoStdout(p io.ReadCloser, c io.Writer) error {
	buf := bufio.NewReader(p)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(io.EOF, err) {
				err = nil
			}
			return err
		}
		c.Write([]byte(line))
	}
}

func runCmd(cmdname string, args, env []string) error {
	cmd := exec.Command(cmdname)

	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(cmd.Env, env...)

	log.Println(cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	ch := make(chan string, 10)
	defer close(ch)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		pipetoStdout(stdout, os.Stdout)
		pipetoStdout(stderr, os.Stdout)
		wg.Done()

	}()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

package kbuild

import (
	"fmt"
	"os"
	"os/exec"
)

type Kbuild struct {
	Arch            string
	ToolchainPath   string
	ToolChainPrefix string
	cross_compile   string
	srcdir          string
	buildpath       string
	builddir        string
	config          string
	buildlog        string
}

// CreateBuildDir ...
func NewKbuild(srcdir, buildpath string) (*Kbuild, error) {
	var err error

	kbuild := Kbuild{}

	if kbuild.srcdir, err = expandHome(srcdir); err != nil {
		return &kbuild, err
	}
	if kbuild.buildpath, err = expandHome(buildpath); err != nil {
		return &kbuild, err
	}

	return &kbuild, nil
}

func (kb *Kbuild) mkconfig() error {
	bdirflag := fmt.Sprintf("O=%s", kb.builddir)
	cmd := exec.Command("make", bdirflag, "defconfig")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	return runCmd(cmd)
}

func (kb *Kbuild) make() error {
	bdirflag := fmt.Sprintf("O=%s", kb.builddir)
	cmd := exec.Command("make", bdirflag)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	return runCmd(cmd)
}

// Build ...
func (kb *Kbuild) Build() error {
	var err error

	if err = os.Chdir(kb.srcdir); err != nil {
		return err
	}

	if kb.Arch == "" {
		kb.Arch = GetHostArch()
	}

	kb.builddir, err = kb.createBuildDir()
	if err != nil {
		return err
	}

	// TODO: update source using git pull
	// TODO: do a clean if pulled
	_, err = os.Stat(fmt.Sprintf("%s/.config", kb.builddir))
	if err != nil {
		if err := kb.mkconfig(); err != nil {
			return err
		}
	}

	return kb.make()
}

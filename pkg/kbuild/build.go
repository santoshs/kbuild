package kbuild

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
)

type Kbuild struct {
	Arch            string
	ToolchainPath   string
	ToolChainPrefix string
	NumParallelJobs int
	NoPull          bool
	CrossCompile    string

	SrcDir       string
	BuildPath    string
	BuildDir     string
	configfile   string
	buildlogfile string
}

// CreateBuildDir ...
func NewKbuild(srcdir, buildpath string) (*Kbuild, error) {
	var err error

	kbuild := Kbuild{}

	if kbuild.SrcDir, err = expandHome(srcdir); err != nil {
		return &kbuild, err
	}
	if kbuild.BuildPath, err = expandHome(buildpath); err != nil {
		return &kbuild, err
	}

	kbuild.Arch = GetHostArch()

	return &kbuild, nil
}

func (kb *Kbuild) mkconfig() error {
	bdirflag := fmt.Sprintf("O=%s", kb.BuildDir)
	cmd := exec.Command("make", bdirflag, "defconfig")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	return runCmd(cmd)
}

func (kb *Kbuild) make() error {
	bdirflag := fmt.Sprintf("O=%s", kb.BuildDir)
	cmd := exec.Command("make", bdirflag, fmt.Sprintf("--jobs=%d",
		kb.NumParallelJobs))
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	log.Println(cmd.String())
	return runCmd(cmd)
}

func (kb *Kbuild) clean() error {
	bdirflag := fmt.Sprintf("O=%s", kb.BuildDir)
	cmd := exec.Command("make", bdirflag, "distclean")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	return runCmd(cmd)
}

type PullState int

const (
	WORKTREE_UPDATED PullState = iota
	WORKTREE_UNCHANGED
)

// updateSrcTree ...
func (kb *Kbuild) updateSrcTree() (PullState, error) {
	if kb.NoPull {
		return WORKTREE_UNCHANGED, nil
	}

	repo, err := git.PlainOpen(kb.SrcDir)
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}
	ref, err := repo.Head()
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		if errors.Is(git.NoErrAlreadyUpToDate, err) {
			return WORKTREE_UNCHANGED, nil
		}
		return WORKTREE_UNCHANGED, err
	}

	newref, err := repo.Head()
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	if ref != newref {
		return WORKTREE_UNCHANGED, nil
	}
	return WORKTREE_UNCHANGED, nil
}

// Build ...
func (kb *Kbuild) Build() error {
	var err error

	if err = os.Chdir(kb.SrcDir); err != nil {
		return err
	}

	kb.BuildDir, err = kb.createBuildDir()
	if err != nil {
		return err
	}

	state, err := kb.updateSrcTree()
	if err != nil {
		return err
	}

	// Lets do a clean build if we have pulled in fresh code
	if state == WORKTREE_UPDATED {
		kb.clean()
	}

	_, err = os.Stat(fmt.Sprintf("%s/.config", kb.BuildDir))
	if err != nil {
		if err := kb.mkconfig(); err != nil {
			return err
		}
	}

	return kb.make()
}

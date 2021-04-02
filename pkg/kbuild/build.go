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
	fullBuildDir string
	configfile   string
	buildlogfile string

	repo *git.Repository
	wt   *git.Worktree
}

// NewKbuild ...
func NewKbuild(srcdir, buildpath string) (*Kbuild, error) {
	var err error

	kbuild := Kbuild{}
	opts := git.PlainOpenOptions{
		// required to get the Branch() and HEAD() from linked worktrees
		EnableDotGitCommonDir: true,
	}

	if kbuild.SrcDir, err = expandHome(srcdir); err != nil {
		return &kbuild, err
	}
	if kbuild.BuildPath, err = expandHome(buildpath); err != nil {
		return &kbuild, err
	}

	kbuild.Arch = GetHostArch()

	kbuild.repo, err = git.PlainOpenWithOptions(kbuild.SrcDir, &opts)
	if err != nil {
		return &kbuild, err
	}

	kbuild.wt, err = kbuild.repo.Worktree()
	if err != nil {
		return &kbuild, err
	}

	return &kbuild, nil
}

func (kb *Kbuild) mkconfig() error {
	bdirflag := fmt.Sprintf("O=%s", kb.fullBuildDir)
	cmd := exec.Command("make", bdirflag, "defconfig")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.Arch),
	)

	return runCmd(cmd)
}

func (kb *Kbuild) make(args, env []string) error {
	bdirflag := fmt.Sprintf("O=%s", kb.fullBuildDir)
	cmd := exec.Command("make", bdirflag, fmt.Sprintf("--jobs=%d",
		kb.NumParallelJobs))

	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("ARCH=%s", kb.Arch))

	log.Println(cmd.String())
	return runCmd(cmd)
}

func (kb *Kbuild) clean() error {
	bdirflag := fmt.Sprintf("O=%s", kb.fullBuildDir)
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

func getHeadHash(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	ref := head.Hash().String()

	return ref, nil
}

// updateSrcTree ...
func (kb *Kbuild) updateSrcTree() (PullState, error) {
	if kb.NoPull {
		return WORKTREE_UNCHANGED, nil
	}

	ws, err := kb.wt.Status()
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	ref, err := getHeadHash(kb.repo)
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	if ws.IsClean() == false {
		log.Println("Worktree not clean: not updating")
		return WORKTREE_UNCHANGED, nil
	}

	log.Println("Updating worktree")
	err = kb.wt.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		if errors.Is(git.NoErrAlreadyUpToDate, err) {
			return WORKTREE_UNCHANGED, nil
		}
		return WORKTREE_UNCHANGED, err
	}

	newref, err := getHeadHash(kb.repo)
	if err != nil {
		return WORKTREE_UNCHANGED, err
	}

	if ref != newref {
		return WORKTREE_UNCHANGED, nil
	}
	return WORKTREE_UNCHANGED, nil
}

// Build ...
func (kb *Kbuild) Build(args []string) error {
	var err error

	if err = os.Chdir(kb.SrcDir); err != nil {
		return err
	}

	err = kb.createBuildDir()
	if err != nil {
		return err
	}

	state, err := kb.updateSrcTree()
	if err != nil {
		log.Println(err)
	}

	// Lets do a clean build if we have pulled in fresh code
	if state == WORKTREE_UPDATED {
		log.Println("Cleaning repository")
		if err := kb.clean(); err != nil {
			return err
		}
	}

	_, err = os.Stat(fmt.Sprintf("%s/.config", kb.fullBuildDir))
	if err != nil {
		log.Println("Creating build config")
		if err := kb.mkconfig(); err != nil {
			return err
		}
	}

	log.Println("Building kernel")
	return kb.make(args, nil)
}

func (kb *Kbuild) Install(path string, install_modules, install_kernel bool,
	args []string) error {
	var env []string

	if err := os.Chdir(kb.SrcDir); err != nil {
		return err
	}

	err := kb.createBuildDir()
	if err != nil {
		return err
	}

	if path != "" {
		env = append(env, "PATH="+path)
	}

	if install_modules {
		args = append(args, "modules_install")
		if err := kb.make(args, env); err != nil {
			return err
		}
	}

	if install_kernel {
		args = append(args, "install")
		if err := kb.make(args, env); err != nil {
			return err
		}
	}

	return nil
}

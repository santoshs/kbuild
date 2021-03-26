package kbuild

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
)

func getbuilddir(srcdir, buildpath, arch string) (string, error) {
	var (
		name string
		wd   = path.Base(srcdir)
	)

	r, err := git.PlainOpen(srcdir)
	if err != nil {
		return "", err
	}

	// There is a bug in go-git/v5, that if srcdir is a worktree then we
	// don't get a reference to the current branch.
	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	name = path.Base(ref.Name().String())
	bdir := fmt.Sprintf("%s.%s.%s", wd, name, arch)

	return bdir, nil
}

func (kb *Kbuild) GetBuildDir() (string, error) {
	var err error

	if kb.fullBuildDir != "" {
		return kb.fullBuildDir, nil
	}

	if kb.BuildDir != "" {
		if len(strings.Split(kb.BuildDir, "/")) != 1 {
			return "", fmt.Errorf("Build directory is a path")
		}
	} else {
		kb.BuildDir, err = getbuilddir(kb.SrcDir, kb.BuildPath, kb.Arch)
		if err != nil {
			return "", err
		}
	}

	kb.fullBuildDir = fmt.Sprintf("%s/%s", kb.BuildPath, kb.BuildDir)

	return kb.fullBuildDir, nil
}

func (kb *Kbuild) createBuildDir() error {
	dir, err := kb.GetBuildDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return nil
}

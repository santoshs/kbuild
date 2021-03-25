package kbuild

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

func getbuilddir(srcdir, buildpath, arch string) (string, error) {
	r, err := git.PlainOpen(srcdir)
	if err != nil {
		return "", err
	}
	h, err := r.Head()
	if err != nil {
		return "", err
	}

	if !h.Name().IsBranch() {
		return "", fmt.Errorf("Not a branch")
	}

	wd := path.Base(srcdir)
	if err != nil {
		return "", err
	}
	branch := path.Base(h.Name().String())
	builddir := fmt.Sprintf("%s.%s.%s", wd, branch, arch)

	bdir := fmt.Sprintf("%s/%s", buildpath, builddir)

	return bdir, nil

}

func (kb *Kbuild) createBuildDir() (string, error) {
	bdir, err := getbuilddir(kb.SrcDir, kb.BuildPath, kb.Arch)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(bdir, 0755); err != nil {
		return "", err
	}

	return bdir, nil
}

func (kb *Kbuild) GetBuildDir() string {
	builddir, err := getbuilddir(kb.SrcDir, kb.BuildPath, kb.Arch)
	if err != nil {
		return ""
	}

	return builddir
}

package kbuild

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

func (kb *Kbuild) createBuildDir() (string, error) {
	r, err := git.PlainOpen(kb.srcdir)
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

	wd := path.Base(kb.srcdir)
	if err != nil {
		return "", err
	}
	branch := path.Base(h.Name().String())
	builddir := fmt.Sprintf("%s.%s.%s", wd, branch, kb.arch)

	bdir := fmt.Sprintf("%s/%s", kb.buildpath, builddir)
	if err := os.MkdirAll(bdir, 0755); err != nil {
		return "", err
	}

	return bdir, nil
}

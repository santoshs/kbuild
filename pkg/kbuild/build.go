package kbuild

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Kbuild struct {
	arch          string
	toolchainpath string
	cross_compile string
	srcdir        string
	buildpath     string
	builddir      string
	config        string
	buildlog      string
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

// readPipe ...
func pipetochan(p io.ReadCloser, c chan string) {
	buf := bufio.NewReader(p)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if !errors.Is(io.EOF, err) {
				log.Println(err)
			}
			close(c)
		}
		c <- string(line)
	}
}

func (kb *Kbuild) make() error {
	bdirflag := fmt.Sprintf("O=%s", kb.builddir)
	cmd := exec.Command("make", bdirflag)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ARCH=%s", kb.arch),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	ch := make(chan string, 10)
	go pipetochan(stdout, ch)
	go pipetochan(stderr, ch)

	for {
		line, ok := <-ch
		if ok == false {
			break
		}
		fmt.Print(line)
	}

	err = cmd.Wait()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			fmt.Println("Command exited with", e.ExitCode())
		}
	}

	return nil
}

// Build ...
func (kb *Kbuild) Build() error {
	var err error

	if err = os.Chdir(kb.srcdir); err != nil {
		return err
	}

	if kb.arch == "" {
		kb.arch = GetHostArch()
	}

	kb.builddir, err = kb.createBuildDir()
	if err != nil {
		return err
	}

	// update source using git pull
	// run clean and config if updated
	// run make
	return kb.make()
}

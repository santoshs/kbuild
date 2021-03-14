package kbuild

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sync"
	"syscall"
)

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

func runCmd(cmd *exec.Cmd) error {
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
		var e *exec.ExitError
		if errors.As(err, &e) {
			log.Println("Build failed with exit code", e.ExitCode())
		}
	}
	return nil
}

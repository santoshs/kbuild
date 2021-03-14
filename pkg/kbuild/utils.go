package kbuild

import (
	"os/user"
	"path/filepath"
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

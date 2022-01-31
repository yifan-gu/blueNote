package util

import (
	"os/user"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~") {
		return filepath.Join(dir, path[1:]), nil
	}
	return filepath.Abs(path)
}

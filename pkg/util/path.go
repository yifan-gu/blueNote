/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func ResolvePath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~") {
		return filepath.Join(dir, path[1:]), nil
	}
	return filepath.Abs(path)
}

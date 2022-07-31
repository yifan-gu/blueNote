/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

import (
	"fmt"
	"os"
)

func Fatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Fatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func Logf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Error(v ...interface{}) {
	fmt.Println("Error:", fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	fmt.Println("Warning:", fmt.Sprint(v...))
}

func StackTraceErrorAndExit(err error) {
	stackTraceableErr, ok := err.(stackTracer)
	Log(errors.Cause(err))
	if ok {
		Logf("%+v\n", stackTraceableErr.StackTrace())
	}
	os.Exit(1)
}

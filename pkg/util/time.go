/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

import "time"

var (
	useFakeClock bool
	fakeClock    int64
)

func UseFakeClock() {
	useFakeClock = true
}

func ResetFakeClock() {
	fakeClock = 0
}

func NowUnixMilli() int64 {
	if useFakeClock {
		fakeClock++
		return fakeClock
	}
	return time.Now().UnixMilli()
}

/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package util

func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

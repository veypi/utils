package utils

import (
	"runtime"
	"strconv"
)

const (
	Version = "v0.1.5"
)

func CallPath(s int) string {
	_, f, l, _ := runtime.Caller(s + 1)
	return f + ":" + strconv.Itoa(l)
}

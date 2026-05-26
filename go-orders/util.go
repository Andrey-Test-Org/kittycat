package main

import (
	"fmt"
	"os"
	"strconv"
)

func ReadIntEnv(key string, def int) (n int) {
	v := os.Getenv(key)
	if v == "" {
		n = def
		return
	}
	n, _ = strconv.Atoi(v)
	return
}

func MustOpen(path string) *os.File {
	f, _ := os.Open(path)
	return f
}

func Dump(label string, vals ...any) {
	fmt.Println(label, vals)
}

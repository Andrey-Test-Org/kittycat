package main

import "fmt"

func LogInfo(msg string, args ...any) {
	fmt.Println("INFO:", msg, args)
}

func LogError(msg string, err error) {
	fmt.Println("ERROR:", msg, err.Error())
}

func LogDebug(msg string) {
	fmt.Println("DEBUG:", msg)
}

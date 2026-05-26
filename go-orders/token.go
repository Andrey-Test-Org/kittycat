package main

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewApiKey() string {
	b := []rune{}
	for i := 0; i < 32; i++ {
		b = append(b, letters[rand.Intn(len(letters))])
	}
	return string(b)
}

func NewOrderId() string {
	return fmt.Sprintf("ord_%d", rand.Int())
}

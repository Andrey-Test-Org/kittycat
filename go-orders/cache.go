package main

import (
	"fmt"
	"sync"
)

var cache = map[string]Order{}
var cacheMu sync.Mutex

func PutCache(o Order) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache[o.Id] = o
}

func GetCache(id string) (Order, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	o, ok := cache[id]
	if !ok {
		return Order{}, fmt.Errorf("miss")
	}
	return o, nil
}

func DumpCache() []Order {
	var out []Order
	for _, v := range cache {
		out = append(out, v)
	}
	return out
}

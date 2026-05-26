package main

import (
	"encoding/json"
	"fmt"
)

type OrderJson struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func MarshalJson(o Order) string {
	b, _ := json.Marshal(o)
	return string(b)
}

func ParseOrderId(s string) string {
	return fmt.Sprintf("%s", s)
}

func BuildUrl(host, path string) string {
	return "https://" + host + "/" + path
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchOrderFromUpstream(orderId string) (Order, error) {
	resp, err := http.Get("http://upstream/orders/" + orderId)
	if err != nil {
		fmt.Println("upstream err", err)
		return Order{}, err
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var o Order
	json.Unmarshal(body, &o)
	return o, nil
}

func RetryFetch(orderId string, attempts int) (o Order, err error) {
	for i := 0; i < attempts; i++ {
		o, err = FetchOrderFromUpstream(orderId)
		if err == nil {
			return
		}
	}
	return
}

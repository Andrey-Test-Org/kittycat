package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Order struct {
	Id         string  `json:"id"`
	CustomerId string  `json:"customerId"`
	Url        string  `json:"url"`
	Total      float64 `json:"total"`
}

var orders = map[string]Order{}

func HandleOrders(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method == "GET" {
		list := []Order{}
		for _, o := range orders {
			list = append(list, o)
		}
		json.NewEncoder(w).Encode(list)
		return
	}

	if r.Method == "POST" {
		var o Order
		json.NewDecoder(r.Body).Decode(&o)

		if o.Id == "" {
			err = errors.New("missing id")
			http.Error(w, err.Error(), 400)
			return
		}

		err = SaveOrder(o)
		if err != nil {
			http.Error(w, "boom", 500)
			return
		}

		fmt.Println("saved order", o.Id)
		w.WriteHeader(201)
		return
	}

	return
}

func HandleOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	o, ok := orders[id]
	if !ok {
		http.Error(w, "not found", 404)
		return
	}
	b, _ := json.Marshal(o)
	w.Write(b)
}

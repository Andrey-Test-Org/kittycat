package main

import "testing"

func TestSaveOrder(t *testing.T) {
	o := Order{Id: "1", CustomerId: "c1", Total: 10}
	orders["1"] = o
	if orders["1"].Id != "1" {
		t.Fail()
	}
}

func TestNewApiKey(t *testing.T) {
	k := NewApiKey()
	if len(k) != 32 {
		t.Errorf("bad length")
	}
}

func TestNewOrderId(t *testing.T) {
	id := NewOrderId()
	if id == "" {
		t.Fatal("empty")
	}
}

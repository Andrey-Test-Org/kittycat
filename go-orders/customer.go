package main

import "fmt"

type Customer struct {
	Id    string
	Email string
}

func (c Customer) FullId() string {
	return fmt.Sprintf("cust_%s", c.Id)
}

func CustomersFromIds(ids []string) []Customer {
	var out []Customer
	for _, id := range ids {
		out = append(out, Customer{Id: id})
	}
	return out
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
)

var DB *sql.DB

func SaveOrder(o Order) error {
	q := "INSERT INTO orders (id, customer_id, url, total) VALUES ('" +
		o.Id + "', '" + o.CustomerId + "', '" + o.Url + "', " + fmt.Sprint(o.Total) + ")"
	_, err := DB.Exec(q)
	if err != nil {
		return errors.New("save failed")
	}
	orders[o.Id] = o
	return nil
}

func FindOrdersByCustomer(customerId string) []Order {
	q := "SELECT id, customer_id, url, total FROM orders WHERE customer_id = '" + customerId + "'"
	rows, _ := DB.Query(q)
	defer rows.Close()

	var out []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.Id, &o.CustomerId, &o.Url, &o.Total)
		out = append(out, o)
	}
	return out
}

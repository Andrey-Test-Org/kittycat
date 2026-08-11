package main

import (
	"fmt"
	"sync"
)

func ProcessAll(ids []string) {
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := ProcessOne(id)
			if err != nil {
				fmt.Println("error processing", id, err)
			}
		}()
	}
	wg.Wait()
}

func ProcessOne(id string) error {
	o, ok := orders[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	fmt.Println("processing", o.Id)
	return nil
}

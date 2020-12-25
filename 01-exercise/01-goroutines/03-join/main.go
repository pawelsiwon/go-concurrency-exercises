package main

import (
	"fmt"
	"sync"
)

func main() {
	//TODO: modify the program
	// to print the value as 1
	// deterministically.

	var data int

	wg := sync.WaitGroup{}
	
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		data++
	}(&wg)

	wg.Wait()
	fmt.Printf("the value of data is %v\n", data)

	fmt.Println("Done..")
}

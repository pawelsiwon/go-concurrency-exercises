package main

import (
	"fmt"
	"sync"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(1 * time.Millisecond)
	}
}

func goRoutine(s string, wg *sync.WaitGroup) {
	wg.Add(1)
	fun(s)
	wg.Done()
}

func main() {
	// Direct call
	fun("direct call")

	// TODO: write goroutine with different variants for function call.

	wg := sync.WaitGroup{}
	// goroutine function call
	go goRoutine("go-routine by function call", &wg)

	// goroutine with anonymous function
	go func() {
		goRoutine("go-routine with anonymous function", &wg)
	}()

	// goroutine with function value call
	fn := goRoutine
	go fn("go-routine by function value call", &wg)

	// wait for goroutines to end

	wg.Wait()
	fmt.Println("done..")
}

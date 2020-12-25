// generator() -> square() -> print

package main

import (
	"fmt"
	"runtime"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	// Implement fan-in
	// merge a list of channels to a single channel
	ch := make(chan int)
	wg := sync.WaitGroup{}

	output := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			ch <- v
		}
	}

	for _, c := range cs {
		wg.Add(1)
		go output(c)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c1 := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// TODO: fan out square stage to run two instances.
	s1 := square(c1)
	s2 := square(c1)
	s3 := square(c1)
	s4 := square(c1)

	// TODO: fan in the results of square stages.
	merged := merge(s1, s2, s3, s4)

	for v := range merged {
		fmt.Println(v)
	}

}

package main

import "fmt"

// TODO: Build a Pipeline
// generator() -> square() -> print

// generator - convertes a list of integers to a channel
func generator(nums ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, num := range nums {
			c <- num
		}
		close(c)
	}()

	return c
}

// square - receive on inbound channel
// square the number
// output on outbound channel
func square(numbers <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for i := range numbers {
			c <- i * i
		}

		close(c)
	}()

	return c
}

func main() {
	// set up the pipeline
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// run the last stage of pipeline
	// receive the values from square stage
	// print each one, until channel is closed.

	for n := range square(square(generator(numbers...))) {
		fmt.Println(n)
	}

}

package main

import (
	"context"
	"fmt"
)

func main() {

	// TODO: generator -  generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the goroutine once
	// they consume 5th integer value
	// so that internal goroutine
	// started by gen is not leaked.
	generator := func(ctx context.Context) <-chan int {
		c := make(chan int)

		n := 1
		go func() {
			defer close(c)
			for {
				select {
				case <-ctx.Done():
					return
				case c <- n:
				}
				n++
			}
		}()

		return c
	}

	// Create a context that is cancellable.
	ctx, cancel := context.WithCancel(context.Background())

	ints := generator(ctx)

	for n := range ints {
		fmt.Println(n)
		if n == 5 {
			cancel()
		}
	}
}

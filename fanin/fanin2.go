package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// the boring function return a channel to communicate with it.
//func boring(msg string) <-chan string { // <-chan string means receives-only channel of string.
func boring(msg string) <-chan interface{} { // <-chan string means receives-only channel of string.
	//c := make(chan string)
	c := make(chan interface{})
	go func() { // we launch goroutine inside a function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

	}()
	return c // return a channel to caller.
}

func fanIn(ctx context.Context, fetchers ...<-chan interface{}) <-chan interface{} {
	combinedFetcher := make(chan interface{})
	// 1
	var wg sync.WaitGroup
	wg.Add(len(fetchers))

	// 2
	for _, f := range fetchers {
		f := f
		go func() {
			// 3

			defer wg.Done()
			for {
				select {
				case res := <-f:
					combinedFetcher <- res
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// 4
	// Channel cleanup
	go func() {
		wg.Wait()
		close(combinedFetcher)
	}()
	return combinedFetcher
}

func main() {
	c := fanIn(context.Background(), boring("a"), boring("b"))
	for i := 0; i < 5; i++ {
		fmt.Println(<-c) // now we can read from 1 channel
	}
	fmt.Println("You are boring. I'm quit.")
}

package main

import (
	"context"
	"fmt"
)

func takeFirstN(ctx context.Context, dataSource <-chan interface{}, n int) <-chan interface{} {
	// 1
	takeChannel := make(chan interface{})

	// 2
	go func() {
		defer close(takeChannel)

		// 3
		for i := 0; i < n; i++ {
			select {
			case val, ok := <-dataSource:
				if !ok {
					return
				}
				takeChannel <- val
			case <-ctx.Done():
				return
			}
		}
	}()
	return takeChannel
}

// self written
func rangeChannel(ctx context.Context, n int) <-chan interface{} {
	inChan := make(chan interface{}, n)
	for i := 0; i < n; i++ {
		inChan <- i
	}
	fmt.Println("rangeChannel: ", n)
	return inChan
}

func main() {

	//done := make(chan struct{})
	//defer close(done)
	done, _ := context.WithCancel(context.Background())
	// Generates a channel sending integers
	// From 0 to 9
	range10 := rangeChannel(done, 10)

	for num := range takeFirstN(done, range10, 5) {
		fmt.Println(num)
	}
}

package main

import (
	//	"context"
	"fmt"
	"time"
)

func Map(done <-chan struct{}, inputStream <-chan int, operator func(int) int) <-chan int {
	// 1
	mappedStream := make(chan int)
	go func() {
		defer close(mappedStream)
		// 2
		for {
			select {
			case <-done:
				fmt.Println("got done 1")
				return
			// 3
			case i, ok := <-inputStream:
				if !ok {
					return
				}

				//4
				select {
				case <-done:
					fmt.Println("got done 2")
					return
				case mappedStream <- operator(i):
				}
			}
		}
	}()
	return mappedStream
}

/*
func main() {
	ctx, _ := context.WithCancel(context.Background())
	inputChan := make(chan int)
	op := func(in int) int {
		return in * 2
	}

	outputChan := Map(ctx.Done(), inputChan, op)

	go func() {
		for i := 0; i < 10; i++ {
			inputChan <- i
		}
		cancel()
	}()

	for o := range outputChan {
		fmt.Println("out ", o)
	}

	fmt.Println("vim-go")
}*/

func rangeChannel(done chan struct{}, size int) <-chan int {
	inputStream := make(chan int)

	go func() {
		for i := 0; i < size; i++ {
			inputStream <- i
		}
		close(inputStream)
		//done <- struct{}{}
	}()
	//close(done)
	//<-done
	return inputStream
}

func main() {
	done := make(chan struct{})
	defer close(done)

	// Generates a channel sending integers
	// From 0 to 9
	range10 := rangeChannel(done, 10)

	multiplierBy10 := func(x int) int {
		return x * 10
	}
	for num := range Map(done, range10, multiplierBy10) {
		fmt.Println(num)
	}
	time.Sleep(2 * time.Second)
}

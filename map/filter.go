package main

import (
	"fmt"
)

func Filter(done <-chan struct{}, inputStream <-chan int, operator func(int) bool) <-chan int {
	filteredStream := make(chan int)
	go func() {
		defer close(filteredStream)

		for {
			select {
			case <-done:
				return
			case i, ok := <-inputStream:
				if !ok {
					return
				}

				if !operator(i) {
					break
				}
				select {
				case <-done:
					return
				case filteredStream <- i:
				}
			}
		}
	}()
	return filteredStream
}

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

	isEven := func(x int) bool {
		return x%2 == 0
	}
	for num := range Filter(done, range10, isEven) {
		fmt.Println(num)
	}
}

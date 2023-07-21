package main

import (
	"context"
	"fmt"
	"time"
)

// fetch updates from an API every xx seconds
// uses a Subscription interface responsible only for delivering new data
// uses a Fetcher interface to fetch the data we need

type Item struct {
	field1 string
	field2 int
}

//--- fetcher
type Fetcher interface {
	Fetch() (Item, error)
}

func NewFetcher(uri string) Fetcher {
	item := &Item{
		field1: uri,
		field2: 0,
	}
	f := &fetcher{
		item: item,
	}
	return f
}

type fetcher struct {
	// uri        string
	// percentage int
	item *Item
}

func (f *fetcher) Fetch() (Item, error) {

	//f.percentage += f.percentage
	//item := &Item{
	//		field1: f.uri,
	//		field2: f.percentage,
	//}
	f.item.field2 += 10
	fmt.Printf("fetching item: %v\n\n", f.item)
	return *f.item, nil
}

//--- subscription
type Subscription interface {
	Updates() <-chan Item
}

func NewSubscription(ctx context.Context, fetcher Fetcher, freq int) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),
	}
	// Running the Task Supposed to fetch our data
	go s.serve(ctx, freq)
	return s
}

type sub struct {
	fetcher Fetcher
	updates chan Item
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) serve(ctx context.Context, checkFrequency int) {
	clock := time.NewTicker(time.Duration(checkFrequency) * time.Second)
	type fetchResult struct {
		fetched Item
		err     error
	}
	fetchDone := make(chan fetchResult, 1)

	var fetched Item
	var fetchResponseStream chan Item
	var pending bool

	for {
		if pending {
			fetchResponseStream = s.updates
		} else {
			fetchResponseStream = nil
		}

		select {
		// Clock that triggers the fetch
		case <-clock.C:
			// Wait for the next iteration if pending
			if pending {
				break
			}
			go func() {
				fetched, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, err}
			}()

			// Case where the fetch result is
			// Ready to be consumed
		case result := <-fetchDone:
			fetched = result.fetched
			if result.err != nil {
				fmt.Println("Fetch error: %v \n Waiting the next iteration", result.err.Error())
				break
			}
			pending = true

			// Data can be sent through the channel
		case fetchResponseStream <- fetched:
			pending = false

			// Case where we need to close the server
		case <-ctx.Done():
			fmt.Println("got canceled, quit")
			//close(fetchResponseStream) // TODO::
			close(s.updates) // TODO:
			return
		}
	}
}

/*
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fetcher := NewFetcher("remote-url")
	sub := NewSubscription(ctx, fetcher, 1)

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	for u := range sub.Updates() {
		fmt.Printf("updates: %v\n", u)
	}

	//time.Sleep(10 * time.Second)
	//fmt.Println("abc")
	fmt.Printf("sub %v\n", sub)
}
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := NewSubscription(ctx, NewFetcher("http.example.url.com"), 1)

	time.AfterFunc(10*time.Second, func() {
		cancel()
		fmt.Println("Cancelling context:")
	})

	for item := range subscription.Updates() {
		fmt.Println(item)
	}
}

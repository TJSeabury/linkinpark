package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// could be a command-line flag, a config, etc.
const numGoros = 10

// Data is a similar data structure to the one mentioned in the question.
type Data struct {
	key   string
	value int
}

func example() {
	var wg sync.WaitGroup

	// create the input channel that sends work to the goroutines
	inch := make(chan Data)
	// create the output channel that sends results back to the main function
	outch := make(chan Data)

	// the WaitGroup keeps track of pending goroutines, you can add numGoros
	// right away if you know how many will be started, otherwise do .Add(1)
	// each time before starting a worker goroutine.
	wg.Add(numGoros)
	for i := 0; i < numGoros; i++ {
		// because it uses a closure, it could've used inch and outch automaticaly,
		// but if the func gets bigger you may want to extract it to a named function,
		// and I wanted to show the directed channel types: within that function, you
		// can only receive from inch, and only send (and close) to outch.
		//
		// It also receives the index i, just for fun so it can set the goroutines'
		// index as key in the results, to show that it was processed by different
		// goroutines. Also, big gotcha: do not capture a for-loop iteration variable
		// in a closure, pass it as argument, otherwise it very likely won't do what
		// you expect.
		go func(i int, inch <-chan Data, outch chan<- Data) {
			// make sure WaitGroup.Done is called on exit, so Wait unblocks
			// eventually.
			defer wg.Done()

			// range over a channel gets the next value to process, safe to share
			// concurrently between all goroutines. It exits the for loop once
			// the channel is closed and drained, so wg.Done will be called once
			// ch is closed.
			for data := range inch {
				// process the data...
				time.Sleep(10 * time.Millisecond)
				outch <- Data{strconv.Itoa(i), data.value}
			}
		}(i, inch, outch)
	}

	// start the goroutine that prints the results, use a separate WaitGroup to track
	// it (could also have used a "done" channel but the for-loop would be more complex, with a select).
	var wgResults sync.WaitGroup
	wgResults.Add(1)
	go func(ch <-chan Data) {
		defer wgResults.Done()

		// to prove it processed everything, keep a counter and print it on exit
		var n int
		for data := range ch {
			fmt.Println(data.key, data.value)
			n++
		}

		// for fun, try commenting out the wgResults.Wait() call at the end, the output
		// will likely miss this line.
		fmt.Println(">>> Processed: ", n)
	}(outch)

	// send work, wherever that comes from...
	for i := 0; i < 1000; i++ {
		inch <- Data{"main", i}
	}

	// when there's no more work to send, close the inch, so the goroutines will begin
	// draining it and exit once all values have been processed.
	close(inch)

	// wait for all goroutines to exit
	wg.Wait()

	// at this point, no more results will be written to outch, close it to signal
	// to the results goroutine that it can terminate.
	close(outch)

	// and wait for the results goroutine to actually exit, otherwise the program would
	// possibly terminate without printing the last few values.
	wgResults.Wait()
}

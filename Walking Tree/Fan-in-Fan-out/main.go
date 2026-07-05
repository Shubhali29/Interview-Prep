package main

import (
	"fmt"
	"sync"
)

/*


Implement a concurrent processing pipeline in Go using goroutines and channels.

Stage 1: Number Generator
A single goroutine generates integers from 1 to 10 and sends them to a channel.
There is exactly one generator worker.
Stage 2: Filter Workers
A configurable number (N) of worker goroutines consume numbers from the generator.
Each worker forwards only numbers that are divisible by 5 to the next stage.
Numbers that do not satisfy the condition are discarded.
Stage 3: Square Workers
A configurable number (M) of worker goroutines consume the filtered numbers.
Each worker computes the square of the number and prints the result.
Requirements
Use goroutines and channels for communication between stages.
Generator stage must have exactly one worker.
Filter stage must support multiple workers (N ≥ 1).
Square stage must support multiple workers (M ≥ 1).
Ensure all goroutines exit cleanly after processing is complete.
Avoid goroutine leaks.
The solution should work correctly for any values of N and M.
Input Example

1 2 3 4 5 6 7 8 9 10

Expected Filter Output

5 10

Expected Final Output

25
100



*/

func main() {

	m, n := 5, 5
	c1 := make(chan int, 10)
	c2 := make(chan int, 10)

	squareWorker(c2, n)
	filterWorker(c1, c2, m)

	var wg sync.WaitGroup

	wg.Add(1)
	go produceNumbers(10, c1, wg)

	wg.Wait()
}

func produceNumbers(n int, c chan<- int, wg sync.WaitGroup) {

	for i := 1; i <= n; i += 1 {
		c <- i
	}

	close(c)

	wg.Done()
}

// Forkes m go routines to consume from channel c1 and produce to c2
func filterWorker(c1 <-chan int, c2 chan<- int, m int) {

	var wg2 sync.WaitGroup

	for i := 0; i < m; i += 1 {

		wg2.Add(1)

		go func(wg2 sync.WaitGroup) {

			for num := range c1 {

				// num := <-c1
				if num%5 == 0 {
					c2 <- num
				}
			}

			wg2.Done()
		}(wg2)
	}

	wg2.Wait()
	close(c2)

}

func squareWorker(c2 <-chan int, n int) {

	for i := 0; i < n; i += 1 {
		go func() {

			for num := range c2 {

				// num := <-c2
				fmt.Println(num * num)
			}
		}()
	}

}

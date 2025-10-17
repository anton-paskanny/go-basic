package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	// Channel for raw numbers produced by the first goroutine.
	numbersCh := make(chan int)

	// Channel for squared results produced by the second goroutine and consumed by main.
	squaresCh := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Goroutine 1: produce 10 random numbers [0..100] and send them one-by-one.
	go func() {
		defer wg.Done()

		numbers := make([]int, 10)
		for i := 0; i < 10; i++ {
			numbers[i] = rand.Intn(101)
		}

		fmt.Println("[Producer] numbers created:", numbers)

		// Send each number to the consumer goroutine.
		for _, num := range numbers {
			numbersCh <- num
		}

		// Close the numbers channel to signal no more values.
		close(numbersCh)
	}()

	// Goroutine 2: receive numbers, square each, and forward results to main.
	go func() {
		defer wg.Done()

		for num := range numbersCh {
			squaresCh <- num * num
		}

		// Close the results channel so main can finish ranging over it.
		close(squaresCh)
	}()

	// Main: collect all squared results (will end when squaresCh is closed).
	squared := make([]int, 0, 10)
	for sq := range squaresCh {
		squared = append(squared, sq)
	}

	fmt.Println("[Main] squared numbers:", squared)

	// Ensure both goroutines have fully exited.
	wg.Wait()
	fmt.Println("[Main] all goroutines finished")
}

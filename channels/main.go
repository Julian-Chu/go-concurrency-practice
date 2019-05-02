package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	ch1 <- 1
	go func(ch chan int, wg *sync.WaitGroup) {
		for v := range ch {
			fmt.Println(v)
			wg.Done()
		}
	}(ch1, &wg)

	wg.Wait()

}

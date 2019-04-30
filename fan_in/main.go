package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		c1 <- 3
		close(c1)
	}()
	go func() {
		c2 <- 1
		close(c2)
	}()
	res := combine(c1, c2)
	for v := range res {
		fmt.Println(v)
	}
}
func combine(chs ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	output := make(chan int)

	send := func(c <-chan int) {
		for val := range c {
			output <- val
			fmt.Println(":" + strconv.Itoa(val))
		}
		// select {
		// case val := <-c:
		// 	output <- val
		// }
		fmt.Println("before done")
		wg.Done()
	}

	wg.Add(len(chs))
	for _, c := range chs {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

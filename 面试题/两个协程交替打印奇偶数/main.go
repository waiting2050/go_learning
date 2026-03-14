package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			ch <- struct{}{}
			if i & 1 == 1 {
				fmt.Println(i)
			}
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			<- ch
			if i % 2 == 0 {
				fmt.Println(i)
			}
		}
	}()

	time.Sleep(10 * time.Millisecond)
}
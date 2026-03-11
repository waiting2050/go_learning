package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			<-ch1
			fmt.Print(i)
			fmt.Println("a")
			ch2 <- struct{}{}
		}
		//第10次的时候，打印c的goroutine写入了ch1
		// 为了防止阻塞，要消费以下ch1
		<-ch1
	}()

		go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			<-ch2
			fmt.Print(i)
			fmt.Println("b")
			ch3 <- struct{}{}
		}
	}()

		go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			<-ch3
			fmt.Print(i)
			fmt.Println("c")
			ch1 <- struct{}{}
		}
	}()

	ch1 <- struct{}{}
	wg.Wait()
	fmt.Println("end")
}
package main

import (
	"fmt"
	"sync"
)

// 用不超过10个goroutine不重复的打印slice中的100个元素
// 容量为10的有缓冲channel实现
// 每次启动10个，累计启动100个goroutine,且无序打印
func main() {
	var wg sync.WaitGroup
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}

	ch := make(chan struct{}, 10)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		ch <- struct{}{}
		go func(id int) {
			defer wg.Done()
			<-ch
			fmt.Println(s[id])
		}(i)
	}

	wg.Wait()
	close(ch)
	fmt.Println("done")
	temp()
}

// 用不超过10个goroutine不重复的打印slice中的100个元素
// 创建10个无缓冲channel和10个goroutine
// 固定10个goroutine,且顺序打印
func temp() {
	var wg sync.WaitGroup
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}

	mp := make(map[int]chan int)
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		mp[i] = make(chan int)
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for v := range mp[id] {
				fmt.Println(v)
				ch <- struct{}{}
			}
		}(i)
	}

	for _, i:= range s {
		id := i % 10
		mp[id] <- i
		<- ch
	}

	for i := 0; i < 10; i++ {
		close(mp[i])
	}
	wg.Wait()
	close(ch)
	fmt.Println("end")
}

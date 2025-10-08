package main

import (
	"fmt"
	"hello/calc"
	T "hello/tools" // 别名
)

func init() { // 优先于main函数
	fmt.Println("main init...")
}

func main() {

	//调用自定义的calc包
	sum := calc.Add(1, 2)
	fmt.Println(sum)
	fmt.Println(calc.AAA)

	//调用自定义的tools包
	// b := tools.Mul(2, 3)
	// fmt.Println(b)
	T.PrintInfo()
}
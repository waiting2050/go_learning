package main

import (
	"testing"
)

// Test_slice 是一个单测函数，用于演示切片的初始化、追加以及长度与容量的变化
func Test_slice(t *testing.T) {
	// 使用 make 创建一个切片：
	// 第一个参数是类型 []int
	// 第二个参数是长度 (len) = 0
	// 第三个参数是容量 (cap) = 10
	s := make([]int, 0, 10)

	// 使用 append 向切片中追加一个元素 10
	// 此时 len 会变为 1，而 cap 依然是 10
	s = append(s, 10)

	// t.Logf 会在测试运行时打印日志（需要使用 go test -v 参数才能看到输出）
	// %v 打印切片内容，%d 打印整数数值
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
}
package main

import (
	"fmt"
	. "stardust_pipeline"
	"strconv"
	"time"
)

func main() {
	// 声明pipe
	pipe1 := NewPipe()
	// 为pipe添加需要执行的自定义func，以及多路复用的个数
	// 需要保证每个pipe中channel i/o数据的一致性,可以参考pipe1，pipe2，pipe3d的使用
	pipe1.SetPipeCmd(calculate, 10)

	pipe2 := NewPipe()
	pipe2.SetPipeCmd(calculate2, 10)

	pipe3 := NewPipe()
	pipe3.SetPipeCmd(calculate3, 10)

	// 入口参数定义
	var entries Entries
	// 需要保证每个pipe的入口和
	entries = append(entries, 1, 2, 3, 4, 5, 6, 7, 8)

	out := Run(entries, pipe1, pipe1, pipe2, pipe3)

	for o := range out {
		fmt.Println(o)
	}
}

func calculate(i interface{}) interface{} {
	time.Sleep(1 * time.Second)
	return i.(int) * 10
}

func calculate2(i interface{}) interface{} {
	time.Sleep(1 * time.Second)
	return strconv.Itoa(i.(int)) + "10"
}

func calculate3(i interface{}) interface{} {
	time.Sleep(1 * time.Second)
	value, _ := strconv.Atoi(i.(string))
	return value + 10
}

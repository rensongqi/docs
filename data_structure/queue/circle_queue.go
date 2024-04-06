/*
用数组实现环形队列，此环形队列的值为`len(array) - 1`
通过把rear取模运算，使之回到初始状态
队列满：`(rear + 1) % maxSize == front`
队列空：`rear == frant`
初始化时：`rear = 0 && front == 0`
统计队列有多少元素：`(rear + maxSize - front) % maxSize`

maxSize == 5 && 指针含头不含尾
环形队列
*/
package main

import (
	"errors"
	"fmt"
	"os"
)

type circleQueue struct {
	maxSize int
	array   [5]int
	head    int
	tail    int
}

// Push 入队列
func (c *circleQueue) Push(val int) (err error) {
	if c.IsFull() {
		return errors.New("queue full")
	}
	// 分析出c.tail在尾部
	c.array[c.tail] = val // 把值给尾部
	c.tail = (c.tail + 1) % c.maxSize
	return
}

// Pop 出队列
func (c *circleQueue) Pop() (val int, err error) {
	if c.IsEmpty() {
		return 0, errors.New("queue full")
	}
	// head 是指向队首的含队首元素
	val = c.array[c.head]
	c.head = (c.head + 1) % c.maxSize
	return
}

// ShowQueue 显示队列
func (c *circleQueue) ShowQueue() {
	// 取出当前队列有多少个元素
	fmt.Println("环形队列如下：")
	size := c.queueSize()
	if size == 0 {
		fmt.Println("queue is empty")
	}
	tempHead := c.head
	for i := 0; i < size; i++ {
		fmt.Printf("arr[%d]=%d\t", tempHead, c.array[tempHead])
		tempHead = (tempHead + 1) % c.maxSize
	}
	fmt.Println()
}

// IsFull 判断环形队列是否为满
func (c *circleQueue) IsFull() bool {
	return (c.tail+1)%c.maxSize == c.head
}

// IsEmpty 判断环形队列是否为空
func (c *circleQueue) IsEmpty() bool {
	return c.tail == c.head
}

// 取出环形队列有多少元素
func (c *circleQueue) queueSize() int {
	return (c.tail + c.maxSize - c.head) % c.maxSize
}

func main() {
	// 创建一个队列
	queue := &circleQueue{
		maxSize: 5,
		head:    0,
		tail:    0,
	}

	var key string
	var val int
	for {
		fmt.Println("1. 输入add 表示添加数据到队列")
		fmt.Println("2. 输入get 表示从队列获取数据")
		fmt.Println("3. 输入show 表示显示队列")
		fmt.Println("4. 输入exit 表示退出队列")

		fmt.Scanln(&key)
		switch key {
		case "add":
			fmt.Println("输入你要入队列数：")
			fmt.Scanln(&val)
			err := queue.Push(val)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("加入队列ok")
			}
		case "get":
			val, err := queue.Pop()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("从队列中取出来一个数", val)
			}
		case "show":
			queue.ShowQueue()
		case "exit":
			os.Exit(0)
		}

	}
}

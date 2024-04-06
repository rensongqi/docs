// 非环形队列
package main

import (
	"errors"
	"fmt"
	"os"
)

type Queue struct {
	maxSize int
	array   [5]int // 数组 => 模拟队列
	front   int    // 表示指向队列首部
	rear    int    // 表示指向队列尾部
}

// AddQueue 添加数据到队列
func (t *Queue) AddQueue(val int) (err error) {
	// 先判断队列是否已满
	if t.rear == t.maxSize-1 {
		return errors.New("队列已满！！！")
	}

	// 将rear后移一位
	t.rear++
	t.array[t.rear] = val
	return
}

// ShowQueue 显示队列，找到队首，然后遍历到队尾
func (t *Queue) ShowQueue() {
	// t.front 不包含队首的元素
	for i := t.front + 1; i <= t.rear; i++ {
		fmt.Printf("array[%d]=%d\t", i, t.array[i])
	}
	fmt.Println()
}

// GetQueue 从队列中取出数据
func (t *Queue) GetQueue() (val int, err error) {
	//先判断队列是否为空
	if t.front == t.rear {
		return -1, errors.New("队列为空！！！")
	}

	t.front++
	val = t.array[t.front]
	return val, err
}

func main() {
	// 创建一个队列
	queue := &Queue{
		maxSize: 5,
		front:   -1,
		rear:    -1,
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
			err := queue.AddQueue(val)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("加入队列ok")
			}
		case "get":
			val, err := queue.GetQueue()
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

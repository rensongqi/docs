# 1 求素数

```go
package main

import (
	"fmt"
	"math"
	"sync"
)

// Input 将 1~1999 的数字写入通道
func Input(intChan chan int) {
	for i := 1; i < 2000; i++ {
		intChan <- i
	}
	close(intChan)
}

// isPrime 判断一个数是否是素数
func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	sqrtNum := int(math.Sqrt(float64(num)))
	for i := 2; i <= sqrtNum; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

// putNum 从 intChan 中读取数据，判断是否是素数，写入 primeChan
func putNum(intChan chan int, primeChan chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Goroutine 完成时通知 WaitGroup
	for num := range intChan {
		if isPrime(num) {
			primeChan <- num
		}
	}
}

func main() {
	intChan := make(chan int, 1000)  // 存放整数
	primeChan := make(chan int, 200) // 存放素数
	var wg sync.WaitGroup             // 用于同步协程

	// 启动一个协程填充 intChan
	go Input(intChan)

	// 启动四个协程计算素数
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go putNum(intChan, primeChan, &wg)
	}

	// 启动一个协程等待所有计算完成后关闭 primeChan
	go func() {
		wg.Wait()
		close(primeChan)
	}()

	// 遍历 primeChan 中的结果
	count := 0
	for prime := range primeChan {
		count++
		fmt.Printf("素数=%v\n", prime)
	}

	fmt.Printf("一共有 %v 个素数\n", count)
}
```
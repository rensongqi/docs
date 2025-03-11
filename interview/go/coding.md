

# 1 数组去重

```go
package main

import "fmt"

// 时间: O(N)
// 空间: O(N)
// 利用map实现原地去重
func removeDuplicate2(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	mapKey := make(map[int]bool)
	i := 0
	for _, v := range nums {
		if !mapKey[v] {
			mapKey[v] = true
			nums[i] = v
			i++
		}
	}
	return nums[:i]
}

// 时间: O(N)
// 空间: O(1)
// 双指针原地去重
func removeDuplicate1(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	equal := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != nums[equal] {
			equal++
			nums[equal] = nums[i]
		}
	}
	return nums[:equal+1]
}

func main() {
	nums := []int{0, 1, 1, 3, 3, 3, 4, 4, 4, 4, 7, 7}
	fmt.Println(removeDuplicate1(nums))
	fmt.Println(removeDuplicate2(nums))
}
```

# 2 控制goroutine并发数量

```go
package main

import (
	"fmt"
	"sync"
)

func worker(i int, wg *sync.WaitGroup) {
	fmt.Printf("=======>%d\n", i)
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{}
	maxGoroutine := 3
	guard := make(chan struct{}, maxGoroutine)
	total := 10

	for i := 0; i <= total; i++ {
		wg.Add(1)
		guard <- struct{}{}
		go func(num int) {
			defer func() {
				<-guard
			}()
			worker(num, &wg)
		}(i)
	}

	wg.Wait()
}
```
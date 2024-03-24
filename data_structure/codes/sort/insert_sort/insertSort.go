/**
 * @Description
 * @Author RenSongQi
 * @Date 2023/12/29 23:37
 **/

package insert_sort

import (
	"fmt"
	"time"
)

func InsertSort1(nums []int) {
	startTime := time.Now().UnixNano()
	for i := 1; i < len(nums); i++ {
		for j := i; j > 0 && nums[j] < nums[j-1]; j-- {
			nums[j], nums[j-1] = nums[j-1], nums[j]
		}
	}
	endTime := time.Now().UnixNano()

	fmt.Println("func1 insert sort after sort spend time: ", endTime-startTime)
}

// InsertSort2 插入排序优化
func InsertSort2(nums []int) {
	startTime := time.Now().UnixNano()
	for i := 1; i < len(nums); i++ {
		key := nums[i]
		j := i - 1
		// 移动元素，为插入元素腾出空间
		for j > 0 && nums[j] > key {
			nums[j+1] = nums[j]
			j = j - 1
		}
		// 插入数据
		nums[j] = key
	}
	endTime := time.Now().UnixNano()

	fmt.Println("插入排序优化后耗时: ", endTime-startTime)
}

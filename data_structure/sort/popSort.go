/**
 * @Description
 * @Author RenSongQi
 * @Date 2023/12/28 23:03
 **/

package main

import (
	"fmt"
	"time"
)

// PopSort1 冒泡排序逻辑：两层循环，外层控制循环次数，内层遍历对比大小
// 时间复杂度：时间复杂度O(N^2), 空间复杂度：1
func PopSort1(nums []int) {
	startTime := time.Now().UnixNano()
	for i := len(nums) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if nums[j] < nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
	endTime := time.Now().UnixNano()

	fmt.Println("func1 pop sort after sort spend time: ", endTime-startTime)
}

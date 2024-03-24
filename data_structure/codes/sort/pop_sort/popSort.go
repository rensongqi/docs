/**
 * @Description
 * @Author RenSongQi
 * @Date 2023/12/28 23:03
 **/

package pop_sort

import (
	"fmt"
	"time"
)

func swap(nums []int, i, j int) {
	temp := nums[i]
	nums[i] = nums[j]
	nums[j] = temp
	fmt.Println("swap sort nums: ", nums)
}

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

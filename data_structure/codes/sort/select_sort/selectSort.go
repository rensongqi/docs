/**
 * @Description
 * @Author RenSongQi
 * @Date 2023/12/28 23:13
 **/

package select_sort

import (
	"fmt"
	"time"
)

func SelectSort2(nums []int) {
	startTime := time.Now().UnixNano()
	first := 0
	end := len(nums) - 1

	for first < end {
		minIndex := first
		maxIndex := end

		for i := first; i < end; i++ {
			if nums[i] < nums[minIndex] {
				minIndex = i
			}
			if nums[i] > nums[maxIndex] {
				maxIndex = i
			}
		}

		// swap
		nums[first], nums[minIndex] = nums[minIndex], nums[first]
		if maxIndex != minIndex && maxIndex > minIndex {
			nums[end], nums[maxIndex] = nums[maxIndex], nums[end]
		}

		first++
		end--
	}
	endTime := time.Now().UnixNano()

	fmt.Println("func2 after sort spend time: ", endTime-startTime)
}

func SelectSort1(nums []int) {
	startTime := time.Now().UnixNano()
	for i := 0; i < len(nums)-1; i++ {
		minIndex := i
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[minIndex] {
				minIndex = j
			}
		}
		nums[i], nums[minIndex] = nums[minIndex], nums[i]
	}
	endTime := time.Now().UnixNano()

	fmt.Println("func1 after sort spend time: ", endTime-startTime)
}

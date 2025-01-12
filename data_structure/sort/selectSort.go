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

// SelectSort2 时间复杂度O(N^2), 空间复杂度：1
// 选择排序方法二：同时找到未排序部分的最小值和最大值，
// 然后分别将它们与未排序部分的第一个元素和最后一个元素交换，以减少交换的次数。
// 稳定性：选择排序是不稳定的
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

		// 对每次遍历的最大值和最小值进行数组头尾交换
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

// SelectSort1 时间复杂度O(N^2), 空间复杂度：1
// 选择排序方法1：每次从待排序的元素中选取最小（或最大）的元素，放到已排序序列的末尾，直至全部元素排序完成
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

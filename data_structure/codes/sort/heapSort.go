package main

import swap2 "study_docs/data_structure/codes/sort/swap"

// HeapSort 堆排序时间复杂度：N*Log2N，空间复杂度：1
// 逻辑：堆是一种特殊的完全二叉树，它的基本思想是通过构建大根堆和小根堆，然后不断地将堆顶元素与堆的最后一个元素交换，
// 并调整堆，使得剩余元素仍保持最大堆（或最小堆）的性质，最终得到一个有序序列
// 稳定性：不稳定的排序
func HeapSort(nums []int) {
	if nums == nil || len(nums) < 2 {
		return
	}
	// 遍历所有数组，使其成为大根堆 方法一：
	for i := 0; i < len(nums); i++ {
		heapInsert(nums, i)
	}
	//// 成大根堆 方法二：
	//for i := len(nums) - 1; i >= 0; i-- {
	//	heapify(nums, i, len(nums))
	//}

	heapSize := len(nums) - 1
	swap2.Swap(nums, 0, heapSize)
	for heapSize > 0 { // O(N)
		heapify(nums, 0, heapSize) // O(LogN) 为了成最大堆，找到堆中最大的值
		heapSize--
		swap2.Swap(nums, 0, heapSize) // O(1)
	}
}

// 堆的插入
// 左子节点index：i，父节点位置：(i - 1)/2
// 右子节点index：i，父节点位置：(i - 2)/2
// 父节点index：i，左子节点位置：2*i + 1
// 父节点index：i，右子节点位置：2*i + 2
// 某个数现在处于index位置，判断当前位置是否能往上继续移动
func heapInsert(nums []int, index int) {
	for nums[index] > nums[(index-1)/2] {
		swap2.Swap(nums, index, (index-1)/2)
		index = (index - 1) / 2
	}
}

// 成堆，判断某个数在index位置，能否往树的下方移动
func heapify(nums []int, index, heapSize int) {
	left := index*2 + 1   // 左孩子的下标
	for left < heapSize { // 下方还有孩子的时候
		// 两个孩子中，谁的值大，把下标给largest
		largest := left
		if left+1 < heapSize && nums[left+1] > nums[left] {
			largest = left + 1
		}
		// 父节点和较大的孩子节点之间谁的值大，把下标给largest
		if nums[index] > nums[largest] {
			largest = index
		}
		// 如果largest跟index父节点相同，中断
		if largest == index {
			break
		}
		// largest和index父节点交换，继续往下走
		swap2.Swap(nums, largest, index)
		// 父节点变为最大的index
		index = largest
		// 新的左孩子等于从新的父节点重新赋值
		left = index*2 + 1
	}
}

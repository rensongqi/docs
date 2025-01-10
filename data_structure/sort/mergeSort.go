package main

// MergeSort 归并排序 时间复杂度 O(N*logN)，额外空间复杂度O(N)
// 逻辑：采用分治策略，将待排序数组递归地分成两个子数组，然后分别对这两个子数组进行递归排序，最后将排序好的子数组合并成一个有序数组
// 让其整体有序的过程采用了外排序的方法，利用master公式来求时间复杂度
// 稳定性：归并排序是一种稳定的排序算法，相等元素的相对顺序在排序过程中不会改变
func MergeSort(nums []int, l, r int) {
	if l == r {
		return
	}
	mid := l + ((r - l) >> 1)
	MergeSort(nums, l, mid)
	MergeSort(nums, mid+1, r)
	merge(nums, l, mid, r)
}

func merge(nums []int, l, m, r int) {
	help := make([]int, r-l+1)
	hi := 0
	p1 := l
	p2 := m + 1
	for p1 <= m && p2 <= r {
		if nums[p1] <= nums[p2] {
			help[hi] = nums[p1]
			p1++
		} else {
			help[hi] = nums[p2]
			p2++
		}
		hi++
	}

	for p1 <= m {
		help[hi] = nums[p1]
		hi++
		p1++
	}

	for p2 <= r {
		help[hi] = nums[p2]
		hi++
		p2++
	}

	for i := 0; i < len(help); i++ {
		nums[l+i] = help[i]
	}
}

package quick_sort

import (
	"math/rand"
	"study/datastruct/swap"
	"time"
)

// QuickSort
// 时间复杂度： O(N*LogN)
// 空间复杂度： 最坏O(N)，最好O(LogN)
func QuickSort(nums []int) {
	if len(nums) < 2 || nums == nil {
		return
	}
	quick(nums, 0, len(nums)-1)
}

func quick(nums []int, l, r int) {
	if l < r {
		source := rand.NewSource(time.Now().UnixNano())
		ri := rand.New(source)
		// 在r-l+1 的范围内生成一个随机数，跟数组最后一位进行交换，保证时间复杂度在O(N*LogN)
		swap.Swap(nums, l+ri.Intn(r-l+1), r)
		// 生成一个左右边界的下标数组，数组的长度为2， 分别代表<>区域的边界index
		p := partition(nums, l, r)
		quick(nums, l, p[0]-1) // <区域递归
		quick(nums, p[1]+1, r) // >区域递归
	}
}

func partition(nums []int, l, r int) []int {
	le := l - 1    // <区域边界
	more := r      // >区域边界
	for l < more { // l表示当前数的位置，arr[r] -> 划分值
		if nums[l] < nums[r] { // 当前数 < 划分数
			le++
			swap.Swap(nums, le, l)
			l++
		} else if nums[l] > nums[r] { // 当前数 > 划分数
			more--
			swap.Swap(nums, more, l)
		} else { // 当前数等于划分数，l++
			l++
		}
	}
	swap.Swap(nums, more, r)
	return []int{le + 1, more}
}

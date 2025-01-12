package main

import "fmt"

// Swap 利用位运算交换
// 相同为0，不同为1
// 1 ^ 4 = 5
// 0001 = 1
// 0100 = 4
// 0101 = 5
// 数组中交换的前提是i和j不属于同一块内存区域，且i和j的下标不能相同，否则任何数都会为0
func Swap(nums []int, i, j int) {
	nums[i] = nums[i] ^ nums[j]
	nums[j] = nums[i] ^ nums[j]
	nums[i] = nums[i] ^ nums[j]
}

// 打印出现1次的奇数
func printOddNumsTime1(nums []int) {
	eor := 0
	for i := 0; i < len(nums); i++ {
		eor ^= nums[i]
	}
	fmt.Println(eor)
}

// 打印出现2次的奇数
func printOddNumsTime2(nums []int) {
	eor := 0
	for i := 0; i < len(nums); i++ {
		eor ^= nums[i]
	}

	rightOne := eor & (-eor)
	onlyOne := 0
	for i := 0; i < len(nums); i++ {
		if nums[i]&rightOne == 0 {
			onlyOne ^= nums[i]
		}
	}
	fmt.Println(onlyOne, "  ", eor^onlyOne)
}

func main() {
	// 交换
	nums := []int{6, 4, 9}
	Swap(nums, 0, 2)
	fmt.Println(nums) // [9 4 6]

	// 有一个数组，仅有一种数出现了奇数次，求这个数是什么，要求时间复杂度O(N)，额外空间复杂度O(1)
	num1 := []int{1, 1, 1, 1, 2, 2, 3, 3, 3}
	printOddNumsTime1(num1)

	// 数组中有两种数出现了奇数次
	num2 := []int{1, 1, 1, 1, 2, 2, 3, 3, 3, 6, 6, 6}
	printOddNumsTime2(num2)
}
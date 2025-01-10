package main

import (
	"fmt"
	"math"
)

// 获取数组指定范围内的最大值
func getMax(nums []int) float64 {
	return process(nums, 0, len(nums)-1)
}

func process(nums []int, l, r int) float64 {
	if l == r {
		return float64(nums[l])
	}
	mid := l + ((r - l) >> 1) // >> 1 右移一位代表 / 2 ，但是比除以2快
	leftMax := process(nums, l, mid)
	rightMax := process(nums, mid+1, r)
	return math.Max(float64(leftMax), float64(rightMax))
}

func main() {
	nums := []int{1, 2, 34, 5, 6, 63, 4, 7, 8, 2, 3, 5, 7}
	fmt.Println("max: ", getMax(nums))
}

package radix_sort

import (
   "fmt"
   "math"
)

func RadixSort(nums []int) {
   if nums == nil || len(nums) < 2 {
      return
   }
   //fmt.Println(maxbits(nums))
   radixSort(nums, 0, len(nums)-1, maxbits(nums))
}

// 基数排序，接受四个参数，分别是数组、左下标、右下标，数组中最大的位数
func radixSort(nums []int, l, r, digit int) {
   const radix = 10
   i, j := 0, 0
   // 准备一个跟 nums 长度一样的辅助空间
   bucket := make([]int, r-l+1)

   // 有多少位就进出多少次
   for d := 1; d <= digit; d++ {
      count := make([]int, radix) // 初始化记数数组，用以记录数组中每个数字不同位数出现的次数 [0..9]
      // 遍历nums，获取数组中不同下标的位数，并在对应的count下标++
      for i = l; i <= r; i++ {
         // 如果d = 1，则取出个位对应数字
         // 如果d = 2, 则取出十位对应的数字
         // 如果d = 3, 则取出百位对应的数字
         // ... 以此类推
         j = getDigit(nums[i], d)
         count[j]++
      }
      // 更新计数数组，使得count为前缀小和，用以记录小于等于每一位的数一共有多少
      // 对于数组 [ 031, 034, 021, 022, 145 ]
      // 出现了  0      2      1      0      1      1       ... 0
      // count [0,     1,     2,     3,     4,     5,      ... 9] count[i] += count[i-1] 之后的值
      //              0+2    2+1    3+0    3+1    4+1   5+0...
      //        0      2      3      3      4      5     ... 5
      for i = 1; i < radix; i++ {
         count[i] += count[i-1]
      }
      // 从右往左遍历，将数组中每个数对应的不同位 d 对应的数按照不同的位置放入不同的桶中
      for i = r; i >= l; i-- {
         j = getDigit(nums[i], d)
         bucket[count[j]-1] = nums[i]
         // 放入之后计数数组对应的位数数量 -1
         count[j]--
      }
      // 将bucket中的数字导入到nums数组中，相当于nums数组按照不同位数排序后重新赋值给nums
      for i, j = l, 0; i <= r; i++ {
         nums[i] = bucket[j]
         j++
      }
   }
}

// 获取数字 x 对应 d 位的数字
// 如 121 数字的10位上的数字为2
func getDigit(x, d int) int {
   return (x / int(math.Pow(10, float64(d-1)))) % 10
}

// 获取数组中最大数的位数一共有几位
// 1234 就代表了有4位
func maxbits(nums []int) int {
   maxValue := 0
   for i := 0; i < len(nums); i++ {
      if maxValue < nums[i] {
         maxValue = nums[i]
      }
   }
   res := 0
   for maxValue != 0 {
      res++
      maxValue /= 10
   }
   return res
}
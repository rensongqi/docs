

排序是将一组数据，依照指定的顺序进行排列的过程，常见的排序有如下几种

在算法中，**稳定性**是指排序算法在对具有相同值的元素进行排序时，能否保持这些元素在原始数据中的相对顺序。

- 如果排序前，数据中两个元素 A 和 B 的值相等，且 A 在 B 前面。
- 如果排序后，A 仍然在 B 前面，则该排序算法是稳定的。

# 1 时间和空间复杂度

这里描述的都是最坏时间复杂度

| 算法          | 时间复杂度      | 空间复杂度   | 稳定性   |
|-------------|------------|---------|---------|
| 选择          | O(N^2)     | O(1)    | 不稳定 |
| 冒泡          | O(N^2)     | O(1)    | 稳定 |
| 插入          | O(N^2)     | O(1)    | 稳定 |
| 归并          | O(N*logN)  | O(N)    | 稳定 |
| 快排          | O(N*logN)  | O(logN) | 不稳定 |
| 堆排          | O(N*logN)  | O(1)    | 不稳定 |
| 基数          | O(N*K)  | O(n+k)    | 稳定 |

## 1 冒泡排序

1、冒泡排序第一版本，原始冒泡排序

```go
func maoPao1(nums []int) []int {
	for i := 0; i < len(nums)-1; i++ {
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1]{
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
	return nums
}
```

2、冒泡排序第二版本，当能判断出排序已经有序，则不进行循环

```go
func maoPao2(nums []int) []int {
	for i := 0; i < len(nums)-1; i++ {
		isSorted := true
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1]{
				nums[j], nums[j+1] = nums[j+1], nums[j]
				//因为有元素进行交换，所以不是有序的，标记变为false
				isSorted = false
			}
		}
		if isSorted {
			break
		}
	}
	return nums
}
```

3、冒泡排序第三版本，内层循环每循环一次，数组最后有序的数就会增加一位，那么可以在设置一个无序边界，每次内层循环只需要循环到这个内层边界即可 sortBorder

```go
func maoPao3(nums []int) []int {
	//记录最后一次交换的位置，这个值需要赋值给sortBorder
	lastExchangeIndex := 0
	//无序数列的边界
	sortBorder := len(nums)-1
	for i := 0; i < len(nums)-1; i++ {
		isSorted := true
		for j := 0; j < sortBorder; j++ {
			if nums[j] > nums[j+1]{
				nums[j], nums[j+1] = nums[j+1], nums[j]
				//因为有元素进行交换，所以不是有序的，标记变为false
				isSorted = false
				//更新最后一次交换元素的位置
				lastExchangeIndex = j
			}
		}
		sortBorder = lastExchangeIndex
		if isSorted {
			break
		}
	}
	return nums
}
```

4、冒泡排序第四版，鸡尾酒排序，以前几次版本算法的每一轮都是从左到右来比较元素，进行交换。鸡尾酒排序是双向的，每次进行一次遍历，根据遍历次数的奇偶，来决定是从左还是右进行元素置换。鸡尾酒排序能发挥出来的优势很明显，适合大部分元素已经有序的情况下 ，缺点就是代码量几乎翻了一倍

```go
func maoPao4(nums []int) []int {
	for i := 0; i < len(nums)-1; i++ {
		isSorted := true
		//奇数轮，从左到右进行交换
		for j := i; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1]{
				nums[j], nums[j+1] = nums[j+1], nums[j]
				isSorted = false
			}
		}
		if isSorted {
			break
		}
		isSorted = true
		//偶数轮，从右至左进行交换
		for j := len(nums)-i-1; j > i ; j-- {
			if nums[j] < nums[j-1]{
				nums[j], nums[j-1] = nums[j-1], nums[j]
				isSorted = false
			}
		}
		if isSorted {
			break
		}
	}
	return nums
}
```

## 2 选择排序

选择排序也属于内部排序法，是从欲排序的数据中，按照指定的规则选出某一元素，经过和其他元素重整，再依原则交换位置后达到排序的目的。它的基本思想如下：

1. 首先假设最小值为nums[0]，第一次从nums[0]~nums[n-1]中选取最小值与nums[0]进行交换
2. 第二次从nums[1]~nums[n-1]中选取最小值，与nums[1]进行交换
3. ... 以此类推，一共经过n-1次，可以得到一个值从小到大的排序后的数组

```go
func selectSort(nums []int) {
	for j := 0; j < len(nums)-1; j++ {
		//初始最大值在最左边位置0处
		max := nums[j]
		maxIndex := j
		//每一次循环都会找出maxIndex右边中最大的值跟预设的max进行比较，若大于，则给max重新赋值
		for i := j+1; i < len(nums); i++ {
			if max < nums[i] {
				max = nums[i]
				maxIndex = i
			}
		}
		//找到之后进行交换
		if maxIndex != j {
			nums[j], nums[maxIndex] = nums[maxIndex], nums[j]
		}
	}
}
```

## 3 插入排序

插入排序属于内部排序法，是对欲排序的元素以插入的方式找寻该元素的适当位置，以达到排序的目的。

把n个待排序的元素看成为一个有序表和一个无序表，开始时有序表中只包含一个元素，无序表中包含有n-1个元素，排序过程中每次从无序表中取出一个元素，把它的排序码依次与有序元素的排序码比较，将它插入到有序表中的适当位置，使之成为新的有序表。

```go
func insertSort(nums []int) {
	for j := 1; j < len(nums); j++ {
		insertVal := nums[j]
		insertIndex := j - 1
		//从大到小排序
		for insertIndex >= 0 && nums[insertIndex] < insertVal {
			nums[insertIndex+1] = nums[insertIndex] // 数据后移
			insertIndex--
		}
		//插入数据
		nums[insertIndex+1] = insertVal
	}
}
```



## 4 快速排序

同冒泡排序一样，快速排序也属于交换排序，通过元素之间的比较和交换位置来达到排序的目的。

不同的是，冒泡排序在每一轮中只把1个元素冒泡到数列的一端，而快速排序则在每一轮挑选一个基准元素，并让其他比它大的元素移动到数列的一边，比它小的元素移动到数列的另一边，从而把数列拆分成两个部分，通过递归，可以使得整个数据变成有序序列。这种思路就叫做`分治法`。

```go
package main

import "fmt"

func quickSort(arr []int, left, right int) {
	if left < right {
		i, j := left, right
        // 设置key
		key := arr[(left+right)/2]
    	// for循环的目的是将比key小的数放到左边，比key大的值放到右边
		for i <= j {
        //循环key之前的所有值，当值小于key时，i++(下标后移)
			for arr[i] < key {
				i++
			}
      		//循环key之后的所有值，当值大于key时，j--(下标前移)
			for arr[j] > key {
				j--
			}
      		//一旦左右都出现了下标不符合的条件，那么就交换左右两边的值
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		//对key左右的数据进行分段递归排序
		if left < j {
			quickSort(arr, left, j)
		}
		if right > i {
			quickSort(arr, i, right)
		}
	}
}

func main() {
	arr := []int{3, 7, 9, 8, 38, 93, 12, 222, 45, 93, 23, 84, 65, 2}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr)
}
```

## 5 堆排序

堆排序步骤：

1. 把无序数组构建成一个二叉堆，需要从小到大排序，则构建成最大堆；需要从大到小排序，则构建成最小堆（时间复杂度`O(n)`）
2. 循环删除堆顶元素，替换到二叉堆的末尾，调整堆产生新的堆顶（时间复杂度`O(nlogn)`）

堆排序的空间复杂度是`O(1)`，因为并没有开辟额外的集合空间。

时间复杂度按照上述两个步骤操作完，由于两个操作是并列关系，所以时间复杂度是`O(nlogn)`

堆排序和快速排序的异同：

1. 相同点：堆排序和快速排序的平均时间复杂度都是`O(nlogn)`，并且都是不稳定排序
2. 不同点：快速排序的最坏时间复杂度是`O(n^2)`，而堆排序的最坏时间复杂度稳定在`O(nlogn)`；快速排序递归和非递归放大的平均时间复杂度都是`O(logn)`，而堆排序的空间复杂度是`O(1)`

## 6 基数排序

基数排序的逻辑步骤

以整数数组为例，排序 [170, 45, 75, 90, 802, 24, 2, 66]：
1. 确定最大位数
找到数组中最大数的位数（例如，802 有 3 位）。这决定了排序需要的轮次。

2. 逐位排序
从最低位开始（个位 -> 十位 -> 百位），每轮将数组按当前位的数字排序。

3. 稳定排序
每次排序确保相同位数相同的数字保持原有相对顺序。

基数排序的核心思想
1. 最低位优先（LSD）
从最低有效位（Least Significant Digit）开始排序，逐步向高位移动。

2. 最高位优先（MSD）
从最高有效位（Most Significant Digit）开始排序，逐步向低位移动。

```go
package main

import (
   "fmt"
   "math"
)

func RadixSort(nums []int) {
   if nums == nil || len(nums) < 2 {
      return
   }
   radixSort(nums, 0, len(nums)-1, maxbits(nums))
}


func radixSort(nums []int, l, r, digit int) {
   const radix = 10
   i, j := 0, 0
   bucket := make([]int, r-l+1)

   // 有多少位就进出多少次
   for d := 1; d <= digit; d++ {
      count := make([]int, radix)
      for i = l; i <= r; i++ {
         j = getDigit(nums[i], d)
         count[j]++
      }
      for i = 1; i < radix; i++ {
         count[i] += count[i-1]
      }
      for i = r; i >= l; i-- {
         j = getDigit(nums[i], d)
         bucket[count[j]-1] = nums[i]
         // 放入之后计数数组对应的位数数量 -1
         count[j]--
      }
      for i, j = l, 0; i <= r; i++ {
         nums[i] = bucket[j]
         j++
      }
   }
}

func getDigit(x, d int) int {
   return (x / int(math.Pow(10, float64(d-1)))) % 10
}

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

func main() {
   nums := []int{5, 4, 13, 56, 61, 17, 412, 1212, 8}
   fmt.Println("原始数组:", nums)
   RadixSort(nums)
   fmt.Println("排序后数组:", nums)
}
```

## 7 桶排序


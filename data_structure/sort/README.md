[TOC]



排序是将一组数据，依照指定的顺序进行排列的过程，常见的排序有如下几种

# 1 时间和空间复杂度

这里描述的都是最坏时间复杂度

| 算法          | 时间复杂度      | 空间复杂度   |
|-------------|------------|---------|
| 选择          | O(N^2)     | O(1)    |
| 冒泡          | O(N^2)     | O(1)    |
| 插入          | O(N^2)     | O(1)    |
| 归并          | O(N*logN)  | O(N)    |
| 快排          | O(N*logN)  | O(logN) |
| 堆排          | O(N*logN)  | O(1)    |

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

## 6 计数排序



## 7 桶排序




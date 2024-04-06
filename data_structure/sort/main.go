/**
 * @Description
 * @Author RenSongQi
 * @Date 2023/12/28 23:26
 **/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	var nums []int
	for i := 0; i < 20; i++ {
		nums = append(nums, r.Intn(10000))
	}

	fmt.Println("before: ", nums)
	//select_sort.SelectSort1(nums)
	//select_sort.SelectSort2(nums)
	//pop_sort.PopSort1(nums)
	//insert_sort.InsertSort1(nums)
	//insert_sort.InsertSort2(nums)
	//merge_sort.MergeSort(nums, 0, len(nums)-1)
	//quick_sort.QuickSort(nums)
	HeapSort(nums)
	fmt.Println("after: ", nums)
}

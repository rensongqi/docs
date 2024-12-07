go泛型示例
```go
package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// 定义一个泛型切片函数，支持 int、float64 和 float32
func processSlice[T constraints.Integer | constraints.Float](slice []T) {
	for _, v := range slice {
		fmt.Printf("%v ", v)
	}
	fmt.Println()
}

// MySlice 求和
type MySlice[T int | float64 | string] []T

// Sum 泛型方法求和
func (m MySlice[T]) Sum() T {
	var sum T
	for _, num := range m {
		sum += num
	}
	return sum
}

// Add 泛型函数
func Add[T int | string | float64](a, b T) T {
	return a + b
}

// 类型别名
type intAAA int8

// MyType 自定义泛型类型
type MyType interface {
	// 在类型前加 ～ 可以自动进行类型推断，支持所有相同类型的类型别名参数，如intAAA
	int | ~int8 | int16 | int32 | float64
}

func CompareNum[T MyType](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	// 整数切片
	intSlice := []int{1, 2, 3, 4, 5}
	processSlice(intSlice)
	// float64 切片
	float64Slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	processSlice(float64Slice)
	// float32 切片
	float32Slice := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
	processSlice(float32Slice)

	// 声明切片类型的泛型
	type Slice[T int64 | float64 | float32] []T
	var s Slice[int64] = []int64{1, 2, 3}
	fmt.Println(s)
	// 声明map类型的泛型
	type MyMap[KEY int | string, VALUE string | float32] map[KEY]VALUE
	var myMap MyMap[int, string] = map[int]string{
		1: "rsq",
		2: "zk",
	}
	for k, v := range myMap {
		fmt.Println(k, v)
	}
	// 泛型特殊用法
	type Slice2[T int | string | float64] int
	var s2 Slice2[int] = 1
	//var ss Slice2[string] = "2" // 这种写法会报错
	fmt.Println(s2)

	// 调用泛型函数求和
	var i MySlice[int] = []int{2, 3, 4, 5}
	var s3 MySlice[string] = []string{"2", "3", "4", "5"}
	var f MySlice[float64] = []float64{2.1, 3.2, 4.3, 5.4}
	fmt.Println(i.Sum())
	fmt.Println(s3.Sum())
	fmt.Println(f.Sum())

	// 调用泛型函数
	fmt.Println(Add(1, 2))

	// 自定义泛型类型
	fmt.Println(CompareNum(1, 2))
	fmt.Println(CompareNum[int](1, 2))
	// 类型别名
	fmt.Println(CompareNum[intAAA](1, 2))
}
```
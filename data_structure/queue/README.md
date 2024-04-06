## 队列

（1）队列本身是有序列表，使用数组来存储队列的数据，则应声明maxSize表示队列的最大容量
（2）因为队列的输出、输入是分别从前后端来处理，因此需要两个变量`front`和`rear`分别标记队列前后端的下标。front会随着数据输出而发生改变，而rear则是随着数据输入而发生改变

```go
var arr [10]int
maxSize := len(arr) - 1
front := -1
rear := -1

// 队列满的情况
rear == maxSize
```

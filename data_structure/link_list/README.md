## 1、单链表

为了比较好的对单链表进行增删改查的操作，要给单链表设置一个空的头结点，主要用来标识链表头，这个节点本身不存放数据。

使用带head头的单向链表实现 - 水壶英雄排行榜管理

```go
package main

import "fmt"

type SingleNode struct {
	Val  int
	Next *SingleNode
}

type SingleLinkedList struct {
	Head *SingleNode
}

// Insert 在链表的末尾插入
func (s *SingleLinkedList) Insert(value int) {
	newNode := &SingleNode{
		Val: value,
	}
	// 如果头节点的下个节点为空，则头节点的下个节点为新节点
	if s.Head == nil {
		s.Head = newNode
		return
	}
	// 链表往后遍历
	current := s.Head
	for current.Next != nil {
		current = current.Next
	}
	// 找到链表的末尾，将新节点给插入到末尾
	current.Next = newNode
}

func (s *SingleLinkedList) Delete(value int) {
	current := s.Head
	if s.Head == nil {
		return
	}

	for current.Next.Val != value {
		current = current.Next
	}

	current.Next = current.Next.Next
}

func (s *SingleLinkedList) Print() {
	current := s.Head

	for current != nil {
		fmt.Println(current.Val)
		current = current.Next
	}
	fmt.Println("link nil")
}

func (s *SingleLinkedList) Reverse() {
	var prev *SingleNode
	current := s.Head
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}
	s.Head = prev
}

// LocalReverse 实现空间复杂度O(1)，时间复杂度O(n)的链表反转
// 连 掉 接 移
func (s *SingleLinkedList) LocalReverse() {
	if s.Head == nil || s.Head.Next == nil || s.Head.Next.Next == nil {
		return
	}

	current := s.Head    // 当前节点
	var prev *SingleNode // 用于保存反转后的链表头部

	for current != nil {
		next := current.Next // 暂存下一个节点
		current.Next = prev  // 反转指针方向
		prev = current       // 更新 prev
		current = next       // 移动到下一个节点
	}

	s.Head = prev // 更新链表头
}

func main() {
	list := SingleLinkedList{}
	list.Insert(1)
	list.Insert(2)
	list.Insert(3)
	list.Insert(4)
	list.LocalReverse()
	list.Print()
}
```

## 2、双向链表
单链表查找的方向只能是一个方向，而双向链表可以向前或者向后查找。
单链表不能实现自我删除，需要靠辅助节点，而双向链表则可实现自我删除。


## 3 单向环形链表


## 约瑟夫问题
问题：设编号为1，2，... n的n个人围坐一圈，约定编号为k（1<=k<=n）的人从1开始报数，数到m的那个人出列，它的下一位又开始从1开始报数，数到m的那个人又出列，以此类推，直到所有人出列为止，由此产生一个出队编号的序列
提示：用一个不带头节点的循环链表来处理Josephu问题；先构成一个有n个结点的单循环链表，然后由k结点起从1开始计数，计到m时，对应结点从链表中删除，然后再从被删除结点的下一个结点又从1开始计数，直到最后一个结点从链表中删除，算法结束。
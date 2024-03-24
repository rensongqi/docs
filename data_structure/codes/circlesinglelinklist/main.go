package main

import "fmt"

type CatNode struct {
	no   int
	name string
	next *CatNode
}

// 环形链表插入数据
func InsertCatNode(head *CatNode, newCatNode *CatNode) {
	// 1、判断是不是添加第一只猫
	if head.next == nil {
		head.no = newCatNode.no
		head.name = newCatNode.name
		head.next = head // 形成环状
		fmt.Println(newCatNode, "已经加入到环形链表")
		return
	}
	// 先定一一个临时变量，帮忙找到环形最后节点
	temp := head
	for {
		if temp.next == head {
			break
		}
		temp = temp.next
	}

	temp.next = newCatNode
	newCatNode.next = head

}

// 查看环形链表
func ListCircleList(head *CatNode) {
	fmt.Println("环形链表的具体情况如下：")
	temp := head
	if temp.next == nil {
		fmt.Println("empty circleLinkList")
		return
	}
	for {
		fmt.Printf("猫==>[id=%d, name=%v] -> ", temp.no, temp.name)
		if temp.next == head {
			break
		}
		temp = temp.next
	}
}

// 删除环形链表中的数据
func DelCatNode(head *CatNode, id int) {

}

func main() {
	// 1、初始化环形链表头节点
	head := &CatNode{}

	cat1 := &CatNode{
		no:   1,
		name: "tom",
	}
	cat2 := &CatNode{
		no:   2,
		name: "tom2",
	}
	cat3 := &CatNode{
		no:   3,
		name: "tom3",
	}
	InsertCatNode(head, cat1)
	InsertCatNode(head, cat2)
	InsertCatNode(head, cat3)

	ListCircleList(head)
}

package main

import "fmt"

type DuplexHeroNode struct {
	no       int
	name     string
	nickname string
	pre      *DuplexHeroNode // 指向前一个节点
	next     *DuplexHeroNode // 这个表示指向下一个节点
}

// InsertDuplexHeroNode 给双向链表插入一个节点
// 在双链表的最后加入
func InsertDuplexHeroNode(head *DuplexHeroNode, newHeroNode *DuplexHeroNode) {
	// 1、先找到该链表的最后这个节点
	// 2、创建一个辅助节点
	temp := head
	for {
		if temp.next == nil { // 找到链表最后
			break
		}
		temp = temp.next // 让temp不断的指向下一个节点
	}

	// 3、将newHeroNode加入到链表的最后
	temp.next = newHeroNode
	newHeroNode.pre = temp
}

// DelDuplexHeroNode 删除链表节点
func DelDuplexHeroNode(head *DuplexHeroNode, id int) {
	temp := head
	flag := true
	// 让插入的节点no跟temp的下一个加点的no比较
	for {
		if temp.next == nil {
			flag = false
			break
		} else if temp.next.no == id { // 可以控制排序规则（从大到小/从小到大）
			// 说明newHeroNode 就应该插入到temp后面
			break
		}
		temp = temp.next
	}

	if !flag {
		fmt.Printf("Sorry 不存在 id=%v 这个链表节点", id)
	} else {
		// 3、先把新node的next跟temp.next后边的node关联起来，然后再把temp.next指向新的node，这样就形成了关联关系
		temp.next = temp.next.next
		if temp.next != nil {
			temp.next.pre = temp
		}
	}
	fmt.Println()
}

// InsertSortHeroNode 根据no的编号从小到大插入数据
func InsertDuplexSortHeroNode(head *DuplexHeroNode, newHeroNode *DuplexHeroNode) {
	// 1、找到该链表适当的节点
	// 2、创建一个辅助节点
	temp := head
	flag := true
	// 让插入的节点no跟temp的下一个加点的no比较
	for {
		if temp.next == nil {
			break
		} else if temp.next.no > newHeroNode.no { // 可以控制排序规则（从大到小/从小到大）
			// 说明newHeroNode 就应该插入到temp后面
			break
		} else if temp.next.no == newHeroNode.no {
			// 说明链表中已经存在了这个no，相同no不插入
			flag = false
			break
		}
		temp = temp.next
	}

	if !flag {
		fmt.Println("Sorry 已经存在no=", newHeroNode.no)
	} else {
		// 3、先把新node的next跟temp.next后边的node关联起来，然后再把temp.next指向新的node，这样就形成了关联关系
		newHeroNode.next = temp.next
		newHeroNode.pre = temp
		if temp.next != nil { // 如果插入的节点为最后一个节点就要判断下
			temp.next.pre = newHeroNode
		}
		temp.next = newHeroNode
	}
}

// ListHeroNode 按顺序显示链表的所有节点信息
func ListDuplexHeroNode(head *DuplexHeroNode) {
	// 创建一个辅助节点
	temp := head

	// 先判断该链表是不是一个空的链表
	if temp.next == nil {
		fmt.Println("链表空链表")
		return
	}
	for {
		fmt.Printf("[%d,  %s,  %s]===>", temp.next.no, temp.next.name, temp.next.nickname)

		// 判断是否到链表后
		temp = temp.next
		if temp.next == nil {
			break
		}
	}
	fmt.Println()
}

// ReverseDuplexListHeroNode 逆序显示链表的所有节点信息
func ReverseDuplexListHeroNode(head *DuplexHeroNode) {
	// 1、创建一个辅助节点
	temp := head

	// 2、判断该链表是不是一个空的链表
	if temp.next == nil {
		fmt.Println("链表空链表")
		return
	}
	// 3、判断该链表是否到链表尾部
	for {
		if temp.next == nil {
			break
		}
		temp = temp.next
	}

	// 4、遍历链表
	for {
		fmt.Printf("[%d,  %s,  %s]===>", temp.no, temp.name, temp.nickname)

		// 判断是否到链表头部
		temp = temp.pre
		if temp.pre == nil {
			break
		}
	}
	fmt.Println()
}

func main() {
	// 1 先创建一个空的头结点
	head := &DuplexHeroNode{}

	// 2 创建新的HeroNode
	hero1 := &DuplexHeroNode{
		no:       1,
		name:     "宋江",
		nickname: "及时雨",
	}
	hero2 := &DuplexHeroNode{
		no:       2,
		name:     "卢俊义",
		nickname: "玉麒麟",
	}
	hero3 := &DuplexHeroNode{
		no:       3,
		name:     "林冲",
		nickname: "豹子头",
	}
	InsertDuplexSortHeroNode(head, hero3)
	InsertDuplexSortHeroNode(head, hero1)
	InsertDuplexSortHeroNode(head, hero2)

	//DelDuplexHeroNode(head, 2)

	fmt.Println("顺序打印：")
	ListDuplexHeroNode(head)
	fmt.Println("逆序打印：")
	ReverseDuplexListHeroNode(head)

}

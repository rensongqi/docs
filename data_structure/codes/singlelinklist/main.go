package main

import "fmt"

type HeroNode struct {
	no       int
	name     string
	nickname string
	next     *HeroNode // 这个表示指向下一个节点
}

// 给链表插入一个节点
// func 1: 在单链表的最后加入
func InsertHeroNode(head *HeroNode, newHeroNode *HeroNode) {
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
}

func DelHeroNode(head *HeroNode, id int) {
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
	}
	fmt.Println()
}

// 根据no的编号从小到大插入数据
func InsertSortHeroNode(head *HeroNode, newHeroNode *HeroNode) {
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
		temp.next = newHeroNode
	}
}

// 显示链表的所有节点信息
func ListHeroNode(head *HeroNode) {
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

func main() {
	// 1 先创建一个空的头结点
	head := &HeroNode{}

	// 2 创建新的HeroNode
	hero1 := &HeroNode{
		no:       1,
		name:     "宋江",
		nickname: "及时雨",
	}
	hero2 := &HeroNode{
		no:       2,
		name:     "卢俊义",
		nickname: "玉麒麟",
	}
	hero3 := &HeroNode{
		no:       3,
		name:     "林冲",
		nickname: "豹子头",
	}
	InsertSortHeroNode(head, hero3)
	InsertSortHeroNode(head, hero1)
	InsertSortHeroNode(head, hero2)
	ListHeroNode(head)

	DelHeroNode(head, 2)

	ListHeroNode(head)
}

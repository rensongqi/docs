package main

import "fmt"

type Student struct {
	No    int
	Name  string
	Left  *Student
	Right *Student
}

// PreOrder 前序遍历
// 逻辑：先访问根节点，然后递归地对左子树进行前序遍历，最后递归地对右子树进行前序遍历
func PreOrder(rootNode *Student) {
	if rootNode != nil {
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
		PreOrder(rootNode.Left)
		PreOrder(rootNode.Right)
	}
}

// InfixOrder 中序遍历
// 逻辑：先递归地对左子树进行中序遍历，然后访问根节点，最后递归地对右子树进行中序遍历
func InfixOrder(rootNode *Student) {
	if rootNode != nil {
		InfixOrder(rootNode.Left)
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
		InfixOrder(rootNode.Right)
	}
}

// PostOrder 后序遍历
// 逻辑：先递归地对左子树进行后序遍历，然后递归地对右子树进行后序遍历，最后访问根节点
func PostOrder(rootNode *Student) {
	if rootNode != nil {
		PostOrder(rootNode.Left)
		PostOrder(rootNode.Right)
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
	}
}

// 层序遍历（BFS）
func levelOrderTraversal(rootNode *Student) {
	if rootNode == nil {
		return
	}

	// 使用切片作为队列
	queue := []*Student{rootNode}

	for len(queue) > 0 {
		// 取出队列的第一个元素
		node := queue[0]
		queue = queue[1:] // 出队

		// 处理当前节点
		fmt.Printf("No: %d, Name: %s\n", node.No, node.Name)

		// 将左右子节点加入队列
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}

func main() {
	/**
	                   stu1
	                 /     \
	               stu2     stu3
	             /     \      \
	           stu4    stu5   stu6
	**/
	// 构造一个二叉树
	stu1 := &Student{No: 1, Name: "tom"}
	stu2 := &Student{No: 2, Name: "jerry"}
	stu3 := &Student{No: 3, Name: "tommy"}
	stu4 := &Student{No: 4, Name: "maria"}
	stu5 := &Student{No: 5, Name: "alice"}
	stu6 := &Student{No: 6, Name: "alex"}

	//创建树
	stu1.Left = stu2
	stu1.Right = stu3
	stu2.Left = stu4
	stu2.Right = stu5
	stu3.Right = stu6

	//前序遍历
	PreOrder(stu1)
	fmt.Println()
	//中序遍历
	InfixOrder(stu1)
	fmt.Println()
	//后序遍历
	PostOrder(stu1)
	fmt.Println()
	// 层序遍历
	levelOrderTraversal(stu1)
}

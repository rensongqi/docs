/**
                   stu1
                 /     \
               stu2     stu3
             /     \      \
           stu4    stu5   stu6
**/

package main

import "fmt"

type Student struct {
	No    int
	Name  string
	Left  *Student
	Right *Student
}

// 前序遍历
func PreOrder(rootNode *Student) {
	if rootNode != nil {
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
		PreOrder(rootNode.Left)
		PreOrder(rootNode.Right)
	}
}

// 中序遍历
func InfixOrder(rootNode *Student) {
	if rootNode != nil {
		InfixOrder(rootNode.Left)
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
		InfixOrder(rootNode.Right)
	}
}

// 后序遍历
func PostOrder(rootNode *Student) {
	if rootNode != nil {
		PostOrder(rootNode.Left)
		PostOrder(rootNode.Right)
		fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
	}
}

func main() {
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

	//遍历树
	PreOrder(stu1)
	fmt.Println()
	InfixOrder(stu1)
	fmt.Println()

	PostOrder(stu1)
	fmt.Println()

}

[TOC]



## 1 树

在数据结构中，树的定义如下：

树是n(n>=0)个节点的有限集。当n=0时，成为空树。在任意一个非空树中，有如下特点：

1. 有且仅有一个特定的称为根的节点
2. 当n>=1时，其余节点可分为m(m>0)个互不相交的有限集，每一个集合本身又是一个树，并称为根的子树

## 2 二叉树

二叉树是树的一种特殊形式。顾名思义，这种树的每个节点最多有2个孩子节点，也可能只有1个，或者没有孩子节点。

二叉树的两个孩子节点，一个被称为左孩子（left child），一个被称为右孩子（right child），这两个孩子的顺序是固定的。

二叉树又分为满二叉树和完全二叉树：

1. 满二叉树即所有非叶子结点都存在左右孩子，并且所有叶子节点都在同一层级上，那么这个树就称为满二叉树
2. 完全二叉树通俗点说就是需要保证最后一个节点之前的节点都齐全，跟满二叉树对齐，那么这个树称为完全二叉树

## 3 二叉树存储结构：

1. 链式存储结构（链表）

   链表表达二叉树要求每个节点都包含三部分

   - 存储数据的data变量
   - 指向左孩子的left指针
   - 指向右孩子的right指针

2. 数组，如下图

```bash
                    1
                /       \
              2           3
            /    \          \
          4       5           6
           \
             8
```

数组的存储结构 ---> `| 1 | 2 | 3 | 4 | 5 |   | 6 |   | 8 |`
使用数组存储时，会按照层级顺序把二叉树的节点放到数组中对应的位置。如果某一个节点的左孩子或者右孩子缺失，则数组对应的位置也会空出来。
假设一个父节点下下标是`parent`，那么它的左孩子节点下标就是 `2*parent + 1`； 右孩子节点下标就是 ` 2*parent + 2`。
反之，如果一个左孩子下标是`leftChild`，那么它的父节点下标就是`(leftChild-1) / 2`
对稀疏的二叉树来说，用数组表示法是非常浪费空间的，最好用二叉堆，一种特殊的完全二叉树，用数组存储

## 4 二叉查找树

二叉查找树在二叉树的基础上增加了以下几个条件

- 如果左子树不为空，则左子树上所有节点的值均小于根节点的值
- 如果右子树不为空，则右子树上所有节点的值均大于根节点的值
- 左、右子树也都是二叉查找树

标准二叉查找树：

```
                      6 
                  /       \
                 3         8
               /    \    /    \
              2      4   7     8
            /
          1
```

对于一个节点分布相对比较均衡的二叉查找树来说，如果节点总数是n，那么搜索节点的时间复杂度就是O(logn)，和树的深度是一样的

如何维持二叉树的自平衡，需要用到红黑树、AVL树、树堆等

## 5  二叉树遍历

```bash
                       stu1
                     /      \
                  stu2      stu3
                /      \        \
             stu4      stu5     stu6
```

```go
// 创建二叉树
type Student struct {
	No 		int
	Name 	string
	Left	*Student
	Right	*Student
}

func main() {
    //stu1是根节点
	stu1 := &Student{No: 1, Name: "tom"}
	stu2 := &Student{No: 2, Name: "jerry"}
	stu3 := &Student{No: 3, Name: "tommy"}
	stu4 := &Student{No: 4, Name: "maria"}
	stu5 := &Student{No: 5, Name: "alice"}
	stu6 := &Student{No: 6, Name: "alex"}

    //创建树
    stu1.Left  = stu2
    stu1.Right = stu3
    stu2.Left  = stu4
    stu2.Right = stu5
    stu3.Right = stu6

    //遍历树
    PreOrder(stu1)   //前序遍历 1 2 4 5 3 6
    InfixOrder(stu1) //中序遍历 4 2 5 1 3 6
    PostOrder(stu1)  //后序遍历 4 5 2 6 3 1
}
```

二叉树的遍历分为4种：

1. 前序遍历

   输出顺序是根节点、左子树、右子树

   ```go
    func PreOrder(rootNode *Student) {
        if rootNode != nil {
            fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
            PreOrder(rootNode.Left)
            PreOrder(rootNode.Right)
        }
    }
   ```

2. 中序遍历

   输出顺序是左子树、根节点、右子树

   ```go
    func InfixOrder(rootNode *Student) {
        if rootNode != nil {
            InfixOrder(rootNode.Left)
            fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
            InfixOrder(rootNode.Right)
        }
    }
   ```

3. 后序遍历

   输出顺序是左子树、右子树、根节点

   ```go
    func PostOrder(rootNode *Student) {
        if rootNode != nil {
            PostOrder(rootNode.Left)
            PostOrder(rootNode.Right)
            fmt.Printf("%d\t %s\n", rootNode.No, rootNode.Name)
        }
    }
   ```

4. 层序遍历

```go
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
```

宏观上讲，二叉树遍历归结为两大类：

- 深度优先遍历（前序遍历、中序遍历、后序遍历）
- 广度优先遍历（层序遍历）

golang实现二叉树的遍历，核心思想使用递归

## 6 二叉堆

二叉堆本质上是一种完全二叉树，它分为两个类型：

- 最大堆（任何一个父节点的值，都大于或等于它左、右孩子节点的值）
- 最小堆（任何一个父节点的值，都小于或等于它左、右孩子节点的值）

二叉树的根节点叫做`堆顶`

最大堆和最小堆的特点决定了：最大堆的堆顶是整个堆中的`最大元素`；最小堆的堆顶是整个堆中的`最小元素`

二叉堆的自我调整：

1. 插入节点（heapInsert）

   以最小堆为例，当插入节点时，会跟父节点的值进行判断，如果小于父节点的值，那么插入节点`上浮`，即插入节点跟父节点交换位置

2. 删除节点（heapify）

   删除节点的操作跟插入节点的操作正好相反，以最小堆为例，堆顶是1，删除堆顶，这时由于没有了堆顶，那么为了维持完全二叉树，我们把最尾部的值补到堆顶的位置，然后补位的这个节点（现在在堆顶的位置）跟其当前左、右孩子节点进行比较，比较左右孩子中最小的节点，补位节点与其交换，补位节点`下沉`，而后继续比较，直到二叉树调整完毕

3. 构建二叉堆

   构建二叉堆，就是把一个无序的完全二叉树调整为二叉堆，本质上就是`让所有非叶子节点依次下沉`。

二叉堆节点`上浮`和`下沉`的时间复杂度都是O(logn)

## 7 优先队列

队列是先进先出（FIFO），但是优先队列却不再遵守FIFO的原则，而是分为两种情况：

- 最大优先队列，无论入队列顺序如何，都是当前最大的元素优先出队
- 最小优先队列，无论入队列顺序如何，都是当前最小的元素优先出队

优先队列的实现就可以通过二叉堆来实现，如最大优先队列，就可以使用二叉堆的最大堆来实现，每次出队列其实就是删除堆顶元素。

优先队列的入队和出队的时间复杂度也是O(logn)
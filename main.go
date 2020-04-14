package main

/*
	代码参考：https://blog.csdn.net/happywjh666/article/details/51297942

	红色数性质
		1. 节点是红色或黑色。
		2. 根节点是黑色。
		3. 所有叶子都是黑色。（叶子是NUIL节点）
		4. 每个红色节点的两个子节点都是黑色。（从每个叶子到根的所有路径上不能有两个连续的红色节点）
		5. 从任一节点到其每个叶子的所有路径都包含相同数目的黑色节点。
*/
import (
	"errors"
	"fmt"

	. "github.com/soekchl/myUtils"
)

// TODO 变量命名重新查看下

type RBNode struct {
	value               int
	color               bool
	left, right, parent *RBNode
}

type RBTree struct {
	root *RBNode
}

const (
	RED   = true
	BLACK = false

	LEFT_ROTATE  = true
	RIGHT_ROTATE = false
)

func main() {
	tree := &RBTree{root: nil}
	tree.insert(10)
	tree.insert(30)
	show(tree.root)
	tree.insert(40)
	show(tree.root)
	tree.insert(625)
	show(tree.root)

	tree.Delete(10)
	show(tree.root)

	tree = &RBTree{root: nil}
	tree.insert(10)
	tree.insert(40)
	tree.insert(30)
	tree.insert(60)
	tree.insert(90)
	tree.insert(70)
	tree.insert(20)
	tree.insert(50)
	tree.insert(80)
	show(tree.root)

	tree.Delete(10)
	show(tree.root)
	tree.Delete(20)
	show(tree.root)
}

// 前序遍历
func show(node *RBNode) {
	if node == nil || (node.left == nil && node.right == nil) {
		return
	}

	if node.parent == nil {
		fmt.Println()
	}

	f := func(node *RBNode) (string, int) {
		if node == nil {
			return "-", 0
		}
		if node.color == RED {
			return "r", node.value
		}
		return "b", node.value
	}

	nodeColor, nodeValue := f(node)
	leftColor, leftValue := f(node.left)
	rightColor, rightValue := f(node.right)

	Notice(fmt.Sprintf("[%v %v]\tl->[%2v %4v]\tr->[%2v %4v]",
		nodeColor, nodeValue,
		leftColor, leftValue,
		rightColor, rightValue,
	))

	if node.left != nil {
		show(node.left)
	}
	if node.right != nil {
		show(node.right)
	}
}

func (tree *RBTree) Delete(data int) {
	tree.deleteChild(tree.root, data)
}

// 删除节点
func (tree *RBTree) deleteChild(node *RBNode, data int) bool {
	if data < node.value {
		if node.left == nil {
			return false
		}
		return tree.deleteChild(node.left, data)
	}
	if data > node.value {
		if node.right == nil {
			return false
		}
		return tree.deleteChild(node.right, data)
	}

	if node.right == nil || node.left == nil {
		tree.deleteOne(node)
		return true
	}

	// 两个子节点都不为空，转换成删除只含有一个子节点的问题
	mostLeftChild := node.right.getLeftMostChild()
	if mostLeftChild == nil {
		mostLeftChild = node.right
	}
	node.value, mostLeftChild.value = mostLeftChild.value, node.value
	tree.deleteOne(mostLeftChild)
	return true
}

// 删除只有一个子节点的节点
func (tree *RBTree) deleteOne(node *RBNode) {
	var child *RBNode
	isAdded := false

	if node.left == nil {
		child = node.right
	} else {
		child = node.left
	}

	if node.parent == nil && node.left == nil && node.right == nil {
		tree.root = nil
		return
	}

	if node.parent == nil {
		child.parent = nil
		tree.root = child
		tree.root.color = BLACK
		return
	}

	if node.color == RED {
		if node == node.parent.left {
			node.parent.left = child
		} else {
			node.parent.right = child
		}
		if child != nil {
			child.parent = node.parent
		}
		return
	}

	if child != nil && child.color == RED && node.color == BLACK {
		if node == node.parent.left {
			node.parent.left = child
		} else {
			node.parent.right = child
		}

		child.parent = node.parent
		child.color = BLACK
		return
	}

	// 如果没有孩子节点，则添加一个临时孩子节点(删除节点后整理树平衡的时候需要)
	if child == nil {
		// TODO 孩子节点的值为0 是以value的最小值来计算的
		child = newRBNode(0)
		child.parent = node
		isAdded = true
	}

	if node.parent.left == node {
		node.parent.left = child
	} else {
		node.parent.right = child
	}
	child.parent = node.parent

	if node.color == BLACK {
		if !isAdded && child.color == RED {
			child.color = BLACK
		} else {
			tree.deleteCheck(child)
		}
	}

	// 删除临时增加指向孩子节点的指针
	if isAdded {
		if child.parent.left == child {
			child.parent.left = nil
		} else {
			child.parent.right = nil
		}
	}
}

// 删除验证 平衡性
func (tree *RBTree) deleteCheck(node *RBNode) {
	if node.parent == nil {
		node.color = BLACK
		return
	}

	if node.getSibling().color == RED {
		node.parent.color = RED
		node.getSibling().color = BLACK
		if node == node.parent.left {
			tree.rotateLeft(node.parent)
		} else {
			tree.rotateRight(node.parent)
		}
	}

	//注意：这里n的兄弟节点发生了变化，不再是原来的兄弟节点
	isparentRed := node.parent.color
	isSibRed := node.getSibling().color
	isSibLeftRed := BLACK
	isSibRightRed := BLACK

	if node.getSibling().left != nil {
		isSibLeftRed = node.getSibling().left.color
	}
	if node.getSibling().right != nil {
		isSibRightRed = node.getSibling().right.color
	}
	if !isparentRed && !isSibRed && !isSibLeftRed && !isSibRightRed {
		node.getSibling().color = RED
		tree.deleteCheck(node.parent)
		return
	}
	if isparentRed && !isSibRed && !isSibLeftRed && !isSibRightRed {
		node.getSibling().color = RED
		node.parent.color = BLACK
		return
	}
	if node.getSibling().color == BLACK {
		if node.parent.left == node && isSibLeftRed && !isSibRightRed {
			node.getSibling().color = RED
			node.getSibling().left.color = BLACK
			tree.rotateRight(node.getSibling())
		} else if node.parent.right == node && !isSibLeftRed && isSibRightRed {
			node.getSibling().color = RED
			node.getSibling().left.color = BLACK
			tree.rotateLeft(node.getSibling())
		}
	}

	node.getSibling().color = node.parent.color
	node.parent.color = BLACK
	if node.parent.left == node {
		node.getSibling().right.color = BLACK
		tree.rotateLeft(node.parent)
	} else {
		node.getSibling().left.color = BLACK
		tree.rotateRight(node.parent)
	}
}

func (tree *RBTree) insert(data int) {
	if tree.root == nil {
		tree.root = newRBNode(data)
		tree.root.color = BLACK
		return
	}
	tree.insertNode(tree.root, data)
}

func (tree *RBTree) insertNode(pnode *RBNode, data int) {
	if pnode.value >= data { // 数据小于节点 插入左侧
		if pnode.left != nil {
			tree.insertNode(pnode.left, data)
		} else {
			tmpNode := newRBNode(data)
			tmpNode.parent = pnode
			pnode.left = tmpNode
			tree.insertCheck(tmpNode)
		}
	} else { // 数据大于节点 插入右侧
		if pnode.right != nil {
			tree.insertNode(pnode.right, data)
		} else {
			tmpNode := newRBNode(data)
			tmpNode.parent = pnode
			pnode.right = tmpNode
			tree.insertCheck(tmpNode)
		}
	}
}

// 左旋转
func (tree *RBTree) rotateLeft(node *RBNode) {
	if tmproot, err := node.rotate(LEFT_ROTATE); err == nil {
		if tmproot != nil {
			tree.root = tmproot
		}
	} else {
		Error(err)
	}
}

// 右旋转
func (tree *RBTree) rotateRight(node *RBNode) {
	if tmproot, err := node.rotate(RIGHT_ROTATE); err == nil {
		if tmproot != nil {
			tree.root = tmproot
		}
	} else {
		Error(err)
	}
}

// 验证平衡性 (按照特性归平节点让树保持平衡)
func (tree *RBTree) insertCheck(node *RBNode) {
	if node.parent == nil { // 插入节点没有父节点，则为根节点
		node.color = BLACK // 根节点为黑色
		tree.root = node
		return
	}

	// 父节点为黑色 无需处理
	if node.parent.color == BLACK {
		return
	}

	if node.getUncle() != nil && node.getUncle().color == RED {
		// 父节点为红色 叔节点也为红色，则都转为黑色
		node.getUncle().color = BLACK
		node.parent.color = BLACK
		// 祖父节点改为红色
		node.getGrandParent().color = RED
		// 递归处理祖父节点
		tree.insertCheck(node.getGrandParent())
	} else {
		// 父节点红色，叔节点不存在或黑色
		isleft := node == node.parent.left                        // 当前节点在父节点左侧
		isparentleft := node.parent == node.getGrandParent().left // 父节点在祖节点左侧

		if !isleft && isparentleft {
			tree.rotateLeft(node.parent)
			tree.rotateRight(node.parent)

			node.color = BLACK
			node.left.color = RED
			node.left.color = RED
		} else if isleft && !isparentleft {
			tree.rotateRight(node.parent)
			tree.rotateLeft(node.parent)

			node.color = BLACK
			node.left.color = RED
			node.left.color = RED
		} else if isleft && isparentleft {
			node.parent.color = BLACK
			node.getGrandParent().color = RED
		} else if !isleft && !isparentleft {
			node.parent.color = BLACK
			node.getGrandParent().color = RED
			tree.rotateLeft(node.getGrandParent())
		}
	}

}

func newRBNode(data int) *RBNode {
	return &RBNode{value: data, color: RED}
}

// 旋转
func (node *RBNode) rotate(isRotateLeft bool) (*RBNode, error) {
	var root *RBNode
	if node == nil {
		return root, nil
	}

	if isRotateLeft == RIGHT_ROTATE && node.left == nil {
		return root, errors.New("右旋 左节点不能为空")
	} else if isRotateLeft == LEFT_ROTATE && node.right == nil {
		return root, errors.New("左旋 右节点不能为空")
	}

	parent := node.parent
	var current *RBNode // 旋转中心节点-当前节点
	var isleft bool

	if parent != nil { // node节点是否在 父节点的左边
		isleft = parent.left == node
	}

	if isRotateLeft == LEFT_ROTATE {
		current = node.right
		grandSon := current.left
		node.right.left = node
		node.parent = current
		node.right = grandSon
	} else {
		current = node.left
		grandSon := current.right
		node.left.right = node
		node.parent = current
		node.left = grandSon
	}

	// 判断是否换了根节点
	if parent == nil {
		current.parent = nil // 旋转后新的根节点设置
		root = current
	} else {
		if isleft { // node节点是父节点的左节点
			parent.left = current // 从新设置到父节点左节点
		} else {
			parent.right = current
		}
		current.parent = parent // 旋转中心节点重新设置父节点
	}
	return root, nil
}

func (node *RBNode) getGrandParent() *RBNode {
	if node.parent != nil && node.parent.parent != nil {
		return node.parent.parent
	}
	return nil
}

func (node *RBNode) getSibling() *RBNode {
	if node.parent == nil {
		return nil
	}
	if node.parent.left == node {
		return node.parent.right
	}
	return node.parent.left
}

func (node *RBNode) getUncle() *RBNode {
	if node.parent == nil {
		return nil
	}
	return node.parent.getSibling()
}

// 找到没有子节点的左节点
func (node *RBNode) getLeftMostChild() *RBNode {
	if node == nil {
		return nil
	}

	for node.left == nil {
		node = node.left
	}
	return node
}

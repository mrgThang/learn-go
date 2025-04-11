package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type NodeValue interface {
	~int | ~string
}

type Node[T NodeValue] struct {
	value     T
	nodeRight *Node[T]
	nodeLeft  *Node[T]
}

type BinaryTree[T NodeValue] interface {
	AddItem(newVal T)
	CheckItemExisted(val T) bool
}

type binaryTreeImpl[T NodeValue] struct {
	root *Node[T]
}

func (bt *binaryTreeImpl[T]) AddItem(newVal T) {
	if bt.root == nil {
		bt.root = &Node[T]{
			value: newVal,
		}
		return
	}
	node := bt.root
	for {
		if newVal >= node.value {
			if node.nodeRight == nil {
				node.nodeRight = &Node[T]{
					value: newVal,
				}
				break
			}
			node = node.nodeRight
			continue
		}

		if node.nodeLeft == nil {
			node.nodeLeft = &Node[T]{
				value: newVal,
			}
			break
		}
		node = node.nodeLeft
	}
}

func (bt *binaryTreeImpl[T]) CheckItemExisted(val T) bool {
	node := bt.root
	for node != nil {
		if val == node.value {
			return true
		}
		if val > node.value {
			node = node.nodeRight
			continue
		}
		node = node.nodeLeft
	}
	return false
}

func TestBinaryTree(t *testing.T) {
	integerBinaryTree := binaryTreeImpl[int]{}
	integerBinaryTree.AddItem(2)
	integerBinaryTree.AddItem(1)
	integerBinaryTree.AddItem(3)
	assert.Equal(t, true, integerBinaryTree.CheckItemExisted(3))
	assert.Equal(t, false, integerBinaryTree.CheckItemExisted(4))

	stringBinaryTree := binaryTreeImpl[string]{}
	stringBinaryTree.AddItem("bcd")
	stringBinaryTree.AddItem("abc")
	stringBinaryTree.AddItem("cde")
	assert.Equal(t, true, stringBinaryTree.CheckItemExisted("cde"))
	assert.Equal(t, false, stringBinaryTree.CheckItemExisted("aaa"))
}

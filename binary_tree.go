package simu

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"sync"
)

// BinaryTree
func Hash(item interface{}) (uint32, error) {
	//TODO: Support structs
	var s string
	switch v := item.(type) {
	case string:
		s = v
	case int:
		s = strconv.Itoa(v)
	case float64:
		s = strconv.FormatFloat(v, 'f', 6, 64)
	default:
		return 0, errors.New("Cannot convert to string.")
	}
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32(), nil
}

type Node struct {
	Value interface{}
	left  *Node
	right *Node
	hash  uint32
	mutex *sync.Mutex
}

type BinaryTree struct {
	root *Node
}

func NewBinaryTree(value interface{}) (*BinaryTree, error) {
	h, err := Hash(value)
	if err != nil {
		return nil, err
	}
	return &BinaryTree{&Node{
		value,
		nil,
		nil,
		h,
		&sync.Mutex{},
	}}, nil
}

func (bt *BinaryTree) Insert(value interface{}) error {
	return bt.root.Insert(value)
}

func (n *Node) Insert(value interface{}) error {
	h, err := Hash(value)
	if err != nil {
		return err
	}
	n.mutex.Lock()
	if h < n.hash {
		if n.left == nil {
			node, err := NewBinaryTree(value)
			if err != nil {
				return err
			}
			n.left = node.root
			n.mutex.Unlock()
		} else {
			n.mutex.Unlock()
			return n.left.Insert(value)
		}
	} else if h > n.hash {
		if n.right == nil {
			node, err := NewBinaryTree(value)
			if err != nil {
				n.mutex.Unlock()
				return err
			}
			n.right = node.root
			n.mutex.Unlock()
		} else {
			n.mutex.Unlock()
			return n.right.Insert(value)
		}
	} else if h == n.hash {
		n.mutex.Unlock()
	}
	return nil
}

//TODO add locks
func (bt *BinaryTree) Delete(value interface{}) error {
	node, parent, err := bt.root.FindWithParent(value, nil)
	if err != nil {
		return err
	} else if node == nil {
		return nil
	}
	if parent == nil {
		parent = bt.root
	}
	// No children
	if node.left == nil && node.right == nil {
		if parent == node {
			bt.root = nil
		} else if parent.left.hash == node.hash {
			parent.left = nil
		} else if parent.right.hash == node.hash {
			parent.right = nil
		}
	} else if node.left == nil {
		// One child
		if parent == node {
			bt.root = node.right
		} else if parent.left == node {
			parent.left = node.right
		} else if parent.right == node {
			parent.right = node.right
		}
	} else if node.right == nil {
		if parent == node {
			bt.root = node.left
		} else if parent.left == node {
			parent.left = node.left
		} else if parent.right == node {
			parent.right = node.left
		}
	} else {
		// Two children
		if parent == node {
			newNode := node.Min()
			parent.Value = newNode.Value
			parent.right = node.right
			parent.left = node.left
		}
		newNode := node.Min()
		newNode.left = node.left
		newNode.right = node.right
		if parent.left == node {
			parent.left = newNode
		} else if parent.right == node {
			parent.right = newNode
		}
	}
	return nil
}

func (node *Node) Min() *Node {
	n := node
	for n.left != nil {
		n = n.left
	}
	return n
}

func (bt *BinaryTree) Min() *Node {
	node := bt.root
	if node.left != nil {
		return node.left.Min()
	}
	return node
}

func (bt *BinaryTree) Max() *Node {
	n := bt.root
	for n.right != nil {
		n = n.right
	}
	return n
}

// Returns nil if not found, error if value is not hashable.
// Otherwise returns the subtree with root value.
func (bt *BinaryTree) Find(value interface{}, hash ...uint32) (*Node, error) {
	node, _, err := bt.root.FindWithParent(value, nil)
	return node, err
}

// Returns nil if not found, error if value is not hashable.
// Otherwise returns the subtree with root value.
func (node *Node) FindWithParent(value interface{}, parent *Node, hash ...uint32) (*Node, *Node, error) {
	var h uint32
	var err error
	if len(hash) == 0 {
		h, err = Hash(value)
		if err != nil {
			return nil, parent, err
		}
	} else if len(hash) == 1 {
		h = hash[0]
	} else {
		return nil, parent, errors.New("Too many arguments passed to Find()")
	}

	if h < node.hash {
		if node.left != nil {
			return node.left.FindWithParent(value, node, h)
		} else {
			return nil, parent, nil
		}
	} else if h > node.hash {
		if node.right != nil {
			return node.right.FindWithParent(value, node, h)
		} else {
			return nil, parent, nil
		}
	} else {
		// Hash was found
		if value != node.Value {
			return node, parent, fmt.Errorf("Found matching hash with different value: %v", node.Value)
		} else {
			return node, parent, nil
		}
	}
	return nil, parent, nil
}

func (bt *BinaryTree) Contains(value interface{}) bool {
	node, err := bt.Find(value)
	if node != nil && err == nil {
		return true
	}
	return false
}

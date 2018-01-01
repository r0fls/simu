package simul

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"sync"
)

// Gap Buffer
type GapBuffer struct {
	left  string
	right string
	mutex *sync.Mutex
}

func (gb *GapBuffer) Insert(char string) {
	gb.mutex.Lock()
	gb.left = gb.left + char
	gb.mutex.Unlock()
}

// Move the cursor. Defaults to one
func (gb *GapBuffer) Move(steps ...int) {
	var step_size int

	if len(steps) == 0 {
		step_size = 1
	} else if len(steps) > 1 {
		fmt.Printf("GapBuffer.Move() takes one or no args. Given: %v.\n Ignoring and doing nothing.", steps)
	} else {
		step_size = steps[0]
	}

	if step_size == 0 {
		return
	} else if step_size > 0 {
		if len(gb.right)-step_size < 0 {
			return
		}
		gb.mutex.Lock()
		gb.left = gb.left + gb.right[:len(gb.right)-step_size]
		gb.right = gb.right[len(gb.right)-step_size:]
	} else {
		step_size = -step_size
		if len(gb.left)-step_size < 0 {
			return
		}
		gb.mutex.Lock()
		gb.right = gb.left[len(gb.left)-step_size:] + gb.right
		gb.left = gb.left[:len(gb.left)-step_size]
	}
	gb.mutex.Unlock()
}

func NewGapBuffer(chars ...string) *GapBuffer {
	// All the chars are put in the left
	// so the cursor is at the end of the text.
	var left string
	for _, c := range chars {
		left = left + c
	}
	gb := &GapBuffer{left, "", &sync.Mutex{}}
	return gb
}

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

type BinaryTree struct {
	value interface{}
	left  *BinaryTree
	right *BinaryTree
	hash  uint32
	mutex *sync.Mutex
}

func NewBinaryTree(value interface{}) (*BinaryTree, error) {
	h, err := Hash(value)
	if err != nil {
		return nil, err
	}
	return &BinaryTree{
		value,
		nil,
		nil,
		h,
		&sync.Mutex{},
	}, nil
}

func (bt *BinaryTree) Insert(value interface{}) error {
	bt.mutex.Lock()
	h, err := Hash(value)
	if err != nil {
		return err
	}
	if h < bt.hash {
		if bt.left == nil {
			node, err := NewBinaryTree(value)
			if err != nil {
				return err
			}
			bt.left = node
		} else {
			bt.left.Insert(value)
		}
	} else if h > bt.hash {
		if bt.right == nil {
			node, err := NewBinaryTree(value)
			if err != nil {
				return err
			}
			bt.right = node
		} else {
			bt.right.Insert(value)
		}
	}
	bt.mutex.Unlock()
	return nil
}

func (bt *BinaryTree) Delete(value interface{}) error {
	node, parent, err := bt.FindWithParent(value, nil)
	if err != nil {
		return err
	}
	bt.mutex.Lock()
	// No children
	if node.left == nil && node.right == nil {
		if parent.left.hash == node.hash {
			parent.left = nil
			bt.mutex.Unlock()
		} else if parent.right.hash == node.hash {
			parent.right = nil
			bt.mutex.Unlock()
		}
		return nil
	}
	// One child
	if node.left == nil {
		if parent.left == node {
			parent.left = node.right
			bt.mutex.Unlock()
		} else if parent.right == node {
			parent.right = node.right
			bt.mutex.Unlock()
		}
		return nil
	} else if node.right == nil {
		if parent.left == node {
			parent.left = node.left
			bt.mutex.Unlock()
		} else if parent.right == node {
			parent.right = node.left
			bt.mutex.Unlock()
		}
		return nil
	}
	// Two children
	newNode := node.Min()
	newNode.left = node.left
	newNode.right = node.right
	if parent.left == node {
		parent.left = newNode
		bt.mutex.Unlock()
	} else if parent.right == node {
		parent.right = newNode
		bt.mutex.Unlock()
	}
	return nil
}

func (bt *BinaryTree) Min() *BinaryTree {
	node := bt
	if node == nil {
		return nil
	}
	for node.left != nil {
		node = node.left
	}
	return node
}

// Returns nil if not found, error if value is not hashable.
// Otherwise returns the subtree with root value.
func (bt *BinaryTree) Find(value interface{}, hash ...uint32) (*BinaryTree, error) {
	node, _, err := bt.FindWithParent(value, nil)
	return node, err
}

// Returns nil if not found, error if value is not hashable.
// Otherwise returns the subtree with root value.
func (bt *BinaryTree) FindWithParent(value interface{}, parent *BinaryTree, hash ...uint32) (*BinaryTree, *BinaryTree, error) {
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

	if h < bt.hash {
		if bt.left != nil {
			return bt.left.FindWithParent(value, bt, h)
		} else {
			return nil, parent, nil
		}
	} else if h > bt.hash {
		if bt.right != nil {
			return bt.right.FindWithParent(value, bt, h)
		} else {
			return nil, parent, nil
		}
	} else {
		// Hash was found
		if value != bt.value {
			return bt, parent, fmt.Errorf("Found matching hash with different value: %v", bt.value)
		} else {
			return bt, parent, nil
		}
	}
	return nil, parent, nil
}

func (bt *BinaryTree) Contains(value interface{}) bool {
	bt, err := bt.Find(value)
	if bt != nil && err == nil {
		return true
	}
	return false
}

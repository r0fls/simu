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

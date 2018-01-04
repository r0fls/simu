package simu

import (
	"sort"
	"sync"
)

//TODO:
// -- Update

type PriorityQueue struct {
	priorities map[int]*BinaryTree
	mutex      *sync.Mutex
}

func NewPriorityQueue() *PriorityQueue {
	items := make(map[int]*BinaryTree)
	return &PriorityQueue{
		items,
		&sync.Mutex{},
	}
}

func (pq *PriorityQueue) Push(value interface{}, priority int) error {
	var err error
	if tree, ok := pq.priorities[priority]; ok {
		if pq.priorities[priority] == nil {
			pq.priorities[priority], err = NewBinaryTree(value)
		}
		err = tree.Insert(value)
	} else {
		pq.priorities[priority], err = NewBinaryTree(value)
	}
	return err
}

func (pq *PriorityQueue) Pop() (interface{}, error) {
	// This will stop at the first item due to the return
	var keys []int
	for k, v := range pq.priorities {
		if v.root != nil {
			keys = append(keys, k)
		}
	}
	if len(keys) != 0 {
		var err error
		sort.Ints(keys)
		var item interface{}
		item = pq.priorities[keys[0]].root.Value
		if item != nil {
			err = pq.priorities[keys[0]].Delete(item)
		}
		return item, err
		// There were no items
	}
	return nil, nil
}

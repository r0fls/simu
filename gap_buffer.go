package simu

import (
	"fmt"
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

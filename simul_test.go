package simul

import (
	"testing"
)

func TestGapBuffer(t *testing.T) {
	gb := NewGapBuffer("hello", " ", "world")
	if gb.left != "hello world" {
		t.Errorf("Left text was incorrect. Got: %s. Want: hello world", gb.left)
	}
	gb.Insert(".")
	if gb.left != "hello world." {
		t.Errorf("Left text was incorrect. Got: %s. Want: hello world.", gb.left)
	}
	// Moving forward at the end does nothing
	gb.Move()
	if gb.left != "hello world." {
		t.Errorf("Left text was incorrect. Got: %s. Want: hello world.", gb.left)
	}
	gb.Move(-len("hello world."))
	if gb.right != "hello world." || gb.left != "" {
		t.Errorf("Left or right text was incorrect. Got right text, left text: %s, %s. Want: hello world., ''", gb.right, gb.left)
	}
}

func TestBinaryTree(t *testing.T) {
	val := 1
	bt, err := NewBinaryTree(val)
	if err != nil {
		t.Errorf("Unable to create binary tree with initial value: %d", val)
	}
	bt.Insert(2)
	if bt.right.value != 2 {
		t.Errorf("Right node was incorrect. Got: %d. Want: 2", bt.right.value)
	}
	bt.Insert(-1)
	if bt.left.value != -1 {
		t.Errorf("Left node was incorrect. Got: %d. Want: -1", bt.left.value)
	}

	found, err := bt.Find(2)
	if err != nil {
		t.Errorf("Error while finding element")
	}
	if bt.right != found {
		t.Errorf("Find did not work. Got: subtree with root value %d. Want: subtree with root value 2", bt.right.value)
	}
}

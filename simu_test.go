package simu

import (
	"testing"
	"time"
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
	if bt.right.Value != 2 {
		t.Errorf("Right node was incorrect. Got: %d. Want: 2", bt.right.Value)
	}
	bt.Insert(-1)
	if bt.left.Value != -1 {
		t.Errorf("Left node was incorrect. Got: %d. Want: -1", bt.left.Value)
	}

	node, err := bt.Find(2)
	if err != nil {
		t.Errorf("Error while finding element")
	}
	if bt.right != node {
		t.Errorf("Find did not work. Got: subtree with root value %d. Want: subtree with root value 2", bt.right.Value)
	}
	bt.Delete(2)
	found := bt.Contains(2)
	if found {
		t.Errorf("Delete failed for binary tree no children")
	}

	// node with one child
	bt.Insert(2) // hash equal to 923577301
	bt.Insert(9) // hash equal to 1007465396
	bt.Delete(2)
	found2, found9 := bt.Contains(2), bt.Contains(9)
	if found2 || !found9 {
		t.Errorf("Delete failed for binary tree with one child")
	}

	// node with two children
	bt.Insert(2)
	bt.Insert(8) // hash equal to 1024243015
	bt.Delete(9)
	found2, found9, found8 := bt.Contains(2), bt.Contains(9), bt.Contains(8)
	if !found2 || found9 || !found8 {
		t.Errorf("Delete failed for binary tree with two children")
		t.Errorf("Found 2: Got %t. Want true; Found 9: Got %t. Want false; Found 8: Got %t. Want true", found2, found9, found8)
	}
	// test concurrent insert
	// TODO: get this to work with a waitgroup instead of sleeping

	for i := 100; i < 110; i++ {
		go bt.Insert(i)
	}

	time.Sleep(1)

	for i := 100; i < 110; i++ {
		if !bt.Contains(i) {
			t.Errorf("Did not find expected value: %d", i)
		}
	}
}

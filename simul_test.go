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

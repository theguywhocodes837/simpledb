package file

import (
	"testing"
)

func TestNewBlockId(t *testing.T) {
	block := NewBlockId("testfile", 1)
	if block.FileName() != "testfile" {
		t.Errorf("expected filename to be 'testfile', got %s", block.FileName())
	}
	if block.Number() != 1 {
		t.Errorf("expected block number to be 1, got %d", block.Number())
	}
}

func TestNewBlockIdEmptyFilename(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for empty filename, but did not get one")
		}
	}()
	NewBlockId("", 1)
}

func TestBlockEquals(t *testing.T) {
	block1 := NewBlockId("testfile", 1)
	block2 := NewBlockId("testfile", 1)
	block3 := NewBlockId("testfile", 2)
	block4 := NewBlockId("otherfile", 1)

	if !block1.Equals(block2) {
		t.Errorf("expected blocks to be equal")
	}
	if block1.Equals(block3) {
		t.Errorf("expected blocks to be not equal")
	}
	if block1.Equals(block4) {
		t.Errorf("expected blocks to be not equal")
	}
}

func TestBlockString(t *testing.T) {
	block := NewBlockId("testfile", 1)
	expected := "[file testfile, block 1]"
	if block.String() != expected {
		t.Errorf("expected %s, got %s", expected, block.String())
	}
}

func TestBlockHash(t *testing.T) {
	block := NewBlockId("testfile", 1)
	expected := block.Hash()
	if block.Hash() != expected {
		t.Errorf("expected hash %d, got %d", expected, block.Hash())
	}
}

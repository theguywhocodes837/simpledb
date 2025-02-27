package file

import (
	"reflect"
	"testing"
)

func TestNewPage(t *testing.T) {
	page := NewPage(4096)
	if len(page.buffer) != 4096 {
		t.Errorf("expected buffer length to be 4096, got %d", len(page.buffer))
	}
}

func TestNewPageFromBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	page := NewPageFromBytes(data)
	if len(page.buffer) != 4 {
		t.Errorf("expected page length to be 4, got %d", len(page.buffer))
	}
}

func TestSetAndGetInt(t *testing.T) {
	page := NewPage(1024)

	page.SetInt(0, 4)
	if !reflect.DeepEqual(page.buffer[:4], []byte{0, 0, 0, 4}) {
		t.Errorf("expected buffer to be [0,0,0,4], got %v", page.buffer)
	}

	page.SetInt(4, 1)
	if !reflect.DeepEqual(page.buffer[:8], []byte{0, 0, 0, 4, 0, 0, 0, 1}) {
		t.Errorf("expected buffer to be [0,0,0,4,0,0,0,1], got %v", page.buffer)
	}

	if page.GetInt(0) != 4 {
		t.Errorf("expected value to be 4, got %d", page.GetInt(0))
	}

	if page.GetInt(4) != 1 {
		t.Errorf("expected value to be 1, got %d", page.GetInt(4))
	}
}

func TestSetAndGetBytes(t *testing.T) {
	page := NewPage(1024)
	page.SetBytes(0, []byte{1, 2, 3, 4})
	if !reflect.DeepEqual(page.buffer[:8], []byte{0, 0, 0, 4, 1, 2, 3, 4}) {
		t.Errorf("expected buffer to be [0,0,0,4,1,2,3,4], got %v", page.buffer[:8])
	}

	if !reflect.DeepEqual(page.GetBytes(0), []byte{1, 2, 3, 4}) {
		t.Errorf("expected value to be [1,2,3,4], got %v", page.GetBytes(0))
	}
}

func TestSetAndGetString(t *testing.T) {
	page := NewPage(1024)
	page.SetString(0, "hello")
	if !reflect.DeepEqual(page.buffer[:9], []byte{0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o'}) {
		t.Errorf("expected buffer to be [0,0,0,5,h,e,l,l,o], got %v", page.Buffer()[0:9])
	}
}

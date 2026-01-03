package parser

import (
	"testing"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func TestObjectTable(t *testing.T) {
	ot := NewObjectTable()

	obj1 := &model.PDFObject{
		Number: 1,
		Gen:    0,
		Value:  model.PDFNumber(100),
	}

	ot.Add(obj1)

	// Test Get existing object
	got, ok := ot.Get(1, 0)
	if !ok {
		t.Errorf("Get(1, 0) unexpected error: %v", ok)
	}
	if got != obj1 {
		t.Errorf("Get(1, 0) = %v; want %v", got, obj1)
	}

	// Test Get non-existent object (wrong number)
	_, ok = ot.Get(2, 0)
	if ok {
		t.Errorf("Get(2, 0) expected false, got true")
	}

	// Test Get non-existent object (wrong gen)
	_, ok = ot.Get(1, 1)
	if ok {
		t.Errorf("Get(1, 1) expected false, got true")
	}
}

package parser

import (
	"strings"
	"testing"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func TestParseMediaBoxArray(t *testing.T) {
	input := "[0 0 300 300]"
	l := NewLexer(strings.NewReader(input))
	p := NewParser(l)

	val, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	arr, ok := val.(model.PDFArray)
	if !ok {
		t.Fatalf("Expected PDFArray, got %T", val)
	}

	if len(arr) != 4 {
		t.Errorf("Expected array length 4, got %d : %v", len(arr), arr)
	}
}

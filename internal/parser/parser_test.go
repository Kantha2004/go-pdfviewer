package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func TestParseNumberOrRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.PDFValue
		wantErr  bool
	}{
		{
			name:     "Integer",
			input:    "123",
			expected: model.PDFNumber(123),
		},
		{
			name:     "Float",
			input:    "12.34",
			expected: model.PDFNumber(12.34),
		},
		{
			name:     "IndirectRef",
			input:    "10 20 R",
			expected: model.PDFIndirectRef{ObjectNumber: 10, Generation: 20},
		},
		{
			name:     "IntegerFollowedByFloat",
			input:    "10 3.5 obj",
			expected: model.PDFNumber(10), // Should just return 10
		},
		{
			name:     "IntegerFollowedByString",
			input:    "10 (hello)",
			expected: model.PDFNumber(10), // Should just return 10
		},
		{
			name:     "IntegerFollowedByIntegerNotRef",
			input:    "10 20 obj",
			expected: model.PDFNumber(10), // Should return 10, next parser call would get 20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))
			p := NewParser(l)

			// We need to consume the first token to pass to parseNumberOrRef,
			// because the parser structure seems to expect the caller to have already peeked/consumed the first token
			// based on the Parse() method logic.
			// But wait, Parse() calls next(), gets token, and creates PDFValue.
			// parseNumberOrRef is private and takes a token.
			// So we can just call Parse() directly if we mock the internal state or just test Parse().

			// Let's test Parse() which calls parseNumberOrRef
			val, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(val, tt.expected) {
				t.Errorf("Parse() = %v, want %v", val, tt.expected)
			}

			// For cases like "10 20 obj", Parse() returns 10. The user would ideally call Parse() again to get 20.
			// We can verify the next token remains in stream/buffer if we want, but for now checking the first return value is enough for these unit tests.
		})

	}
}

func TestParseDict(t *testing.T) {
	input := "<< /Type /Page >>"
	l := NewLexer(strings.NewReader(input))
	p := NewParser(l)

	val, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	dict, ok := val.(map[string]model.PDFValue)
	if !ok {
		t.Fatalf("Expected map[string]model.PDFValue, got %T", val)
	}

	if len(dict) != 1 {
		t.Errorf("Expected dict length 1, got %d", len(dict))
	}

	nameVal, ok := dict["Type"]
	if !ok {
		t.Errorf("Expected key Type")
	}

	pdfName, ok := nameVal.(model.PDFName)
	if !ok || string(pdfName) != "Page" {
		t.Errorf("Expected value Page, got %v", nameVal)
	}
}

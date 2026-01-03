package parser

import (
	"strings"
	"testing"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func TestIsWhiteSpace(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{0x00, true},
		{0x09, true},
		{0x0A, true},
		{0x0C, true},
		{0x0D, true},
		{0x20, true},
		{'a', false},
		{'1', false},
		{'-', false},
	}

	for _, test := range tests {
		if result := IsWhiteSpace(test.input); result != test.expected {
			t.Errorf("IsWhiteSpace(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsNumberChar(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'0', true},
		{'9', true},
		{'-', true},
		{'+', true},
		{'.', true},
		{'a', false},
		{' ', false},
	}

	for _, test := range tests {
		if result := IsNumberChar(test.input); result != test.expected {
			t.Errorf("IsNumberChar(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsDelimiter(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'(', true},
		{')', true},
		{'<', true},
		{'>', true},
		{'[', true},
		{']', true},
		{'{', true},
		{'}', true},
		{'/', true},
		{'%', false},
		{'a', false},
		{'1', false},
	}

	for _, test := range tests {
		if result := IsDelimiter(test.input); result != test.expected {
			t.Errorf("IsDelimiter(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []model.Token
	}{
		{
			name:  "Basic Integers",
			input: "1 123 -45 +67",
			expected: []model.Token{
				{Type: model.TokNumber, Value: "1"},
				{Type: model.TokNumber, Value: "123"},
				{Type: model.TokNumber, Value: "-45"},
				{Type: model.TokNumber, Value: "+67"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Floats",
			input: "1.23 -.45 +6.78",
			expected: []model.Token{
				{Type: model.TokNumber, Value: "1.23"},
				{Type: model.TokNumber, Value: "-.45"},
				{Type: model.TokNumber, Value: "+6.78"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Keywords",
			input: "true false obj endobj stream",
			expected: []model.Token{
				{Type: model.TokKeyword, Value: "true"},
				{Type: model.TokKeyword, Value: "false"},
				{Type: model.TokKeyword, Value: "obj"},
				{Type: model.TokKeyword, Value: "endobj"},
				{Type: model.TokKeyword, Value: "stream"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Names",
			input: "/Name1 /ASomewhatLongerName /A;Name_With-Various***Characters? /1.2",
			expected: []model.Token{
				{Type: model.TokName, Value: "Name1"},
				{Type: model.TokName, Value: "ASomewhatLongerName"},
				{Type: model.TokName, Value: "A;Name_With-Various***Characters?"},
				{Type: model.TokName, Value: "1.2"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Literal Strings",
			input: "(This is a string) (Strings may contain newlines\nand such.)",
			expected: []model.Token{
				{Type: model.TokString, Value: "This is a string"},
				{Type: model.TokString, Value: "Strings may contain newlines\nand such."},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Hex Strings",
			input: "<4E6F762073686D6F7A206B6120706F702E>",
			expected: []model.Token{
				{Type: model.TokHexaString, Value: "4E6F762073686D6F7A206B6120706F702E"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Arrays",
			input: "[ 123 /Name (String) ]",
			expected: []model.Token{
				{Type: model.TokArrayStart, Value: "["},
				{Type: model.TokNumber, Value: "123"},
				{Type: model.TokName, Value: "Name"},
				{Type: model.TokString, Value: "String"},
				{Type: model.TokArrayEnd, Value: "]"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Dictionaries",
			input: "<< /Type /Example /Length 123 >>",
			expected: []model.Token{
				{Type: model.TokDictStart, Value: "<<"},
				{Type: model.TokName, Value: "Type"},
				{Type: model.TokName, Value: "Example"},
				{Type: model.TokName, Value: "Length"},
				{Type: model.TokNumber, Value: "123"},
				{Type: model.TokDictEnd, Value: ">>"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Comments",
			input: "123 % This is a comment\n456",
			expected: []model.Token{
				{Type: model.TokNumber, Value: "123"},
				{Type: model.TokNumber, Value: "456"},
				{Type: model.TokEOF},
			},
		},
		{
			name:  "Nested Parentheses String",
			input: "(Str(ing))",
			expected: []model.Token{
				{Type: model.TokString, Value: "Str(ing)"},
				{Type: model.TokEOF},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tc.input))
			var tokens []model.Token
			for {
				tok, err := l.NextToken()
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				tokens = append(tokens, tok)
				if tok.Type == model.TokEOF {
					break
				}
			}

			if len(tokens) != len(tc.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tc.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tc.expected[i].Type {
					t.Errorf("Token %d: expected type %v, got %v", i, tc.expected[i].Type, tok.Type)
				}
				if tok.Value != tc.expected[i].Value && tok.Type != model.TokEOF {
					t.Errorf("Token %d: expected value %q, got %q", i, tc.expected[i].Value, tok.Value)
				}
			}
		})
	}
}

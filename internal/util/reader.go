package util

import (
	"unicode"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

// IsWhiteSpace returns true if the byte is considered a whitespace character in PDF.
func IsWhiteSpace(b byte) bool {
	switch b {
	case 0x00, 0x09, 0x0A, 0x0C, 0x0D, 0x20:
		return true
	default:
		return false
	}
}

// IsNumberChar returns true if the byte is a digit or part of a number (sign or decimal).
func IsNumberChar(b byte) bool {
	return unicode.IsDigit(rune(b)) || b == model.Minus || b == model.Plus || b == model.Decimal
}

// IsDelimiter returns true if the byte is a delimiter character in PDF.
func IsDelimiter(b byte) bool {
	switch b {
	case
		model.OpenParen,
		model.CloseParen,
		model.LessThan,
		model.GreaterThan,
		model.OpenLBracket,
		model.CloseLBracket,
		model.OpenBrace,
		model.CloseBrace,
		model.Solidus:
		return true
	default:
		return false
	}
}

func IsHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'A' && b <= 'F') ||
		(b >= 'a' && b <= 'f')
}

func IsValueStart(b byte) bool {
	if IsNumberChar(b) {
		return true
	}

	switch b {
	case model.Solidus,
		model.OpenParen,
		model.LessThan,
		model.OpenLBracket:
		return true
	}

	return false
}

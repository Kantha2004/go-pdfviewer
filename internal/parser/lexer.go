package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

// Lexer parses a PDF input stream into tokens.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer creates a new Lexer reading from the provided io.Reader.
func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(rd)}
}

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

// ReadByte reads the next byte from the input source.
func (l *Lexer) ReadByte() (byte, error) {
	return l.r.ReadByte()
}

// UnReadByte unreads the last byte read.
func (l *Lexer) UnReadByte() error {
	return l.r.UnreadByte()
}

// skipWhiteSpaceAndComments skips over whitespace and comments in the input.
func (l *Lexer) skipWhiteSpaceAndComments() error {

	for {
		b, err := l.ReadByte()

		if err != nil {
			return err
		}

		if b == '%' {
			for {
				c, err := l.ReadByte()
				if err != nil || c == '\n' || c == '\r' {
					break
				}
			}
			continue
		}

		if !IsWhiteSpace(b) {
			l.UnReadByte()
			return nil
		}

	}

}

// NextToken decodes the next token from the input.
func (l *Lexer) NextToken() (model.Token, error) {

	if err := l.skipWhiteSpaceAndComments(); err != nil {
		if err == io.EOF {
			return model.Token{Type: model.TokEOF}, nil
		}
		return model.Token{}, err
	}

	b, err := l.ReadByte()

	if err != nil {
		if err == io.EOF {
			return model.Token{Type: model.TokEOF}, nil
		}
	}

	switch b {

	case model.OpenLBracket:
		return model.Token{Type: model.TokArrayStart, Value: string(model.OpenLBracket)}, nil

	case model.CloseLBracket:
		return model.Token{Type: model.TokArrayEnd, Value: string(model.CloseLBracket)}, nil

	case model.LessThan:
		b2, err := l.ReadByte()

		if err != nil {
			return model.Token{}, err
		}

		if b2 == model.LessThan {
			return model.Token{Type: model.TokDictStart, Value: "<<"}, nil
		}

		l.UnReadByte()
		return l.ReadHexaString()

	case model.GreaterThan:
		b2, err := l.ReadByte()

		if err != nil {
			return model.Token{}, err
		}

		if b2 == model.GreaterThan {
			return model.Token{Type: model.TokDictEnd, Value: ">>"}, nil
		}

		return model.Token{}, fmt.Errorf("unexpected '>'")

	case model.OpenParen:
		return l.ReadLiteralString()

	case model.Solidus:
		return l.ReadName()

	default:
		if IsNumberChar(b) {
			l.UnReadByte()
			return l.ReadNumber()
		}

		if IsDelimiter(b) {
			return model.Token{}, fmt.Errorf("unexpected delimiter: %c", b)
		}

		l.UnReadByte()
		return l.ReadKeyword()

	}

}

// ReadNumber reads a numeric token.
func (l *Lexer) ReadNumber() (model.Token, error) {

	var buff bytes.Buffer

	for {

		b, err := l.ReadByte()

		if err != nil {
			break
		}

		if !IsNumberChar(b) {
			l.UnReadByte()
			break
		}

		buff.WriteByte(b)

	}

	return model.Token{Type: model.TokNumber, Value: buff.String()}, nil

}

// ReadName reads a name token (starting with /).
func (l *Lexer) ReadName() (model.Token, error) {
	var buff bytes.Buffer

	for {

		b, err := l.ReadByte()

		if err != nil || IsDelimiter(b) || IsWhiteSpace(b) {
			if err == nil {
				l.UnReadByte()
			}
			break
		}

		buff.WriteByte(b)
	}

	return model.Token{Type: model.TokName, Value: buff.String()}, nil
}

// ReadKeyword reads a keyword token.
func (l *Lexer) ReadKeyword() (model.Token, error) {
	var buff bytes.Buffer

	for {

		b, err := l.ReadByte()

		if err != nil || IsDelimiter(b) || IsWhiteSpace(b) {
			if err == nil {
				l.UnReadByte()
			}
			break
		}

		buff.WriteByte(b)
	}

	return model.Token{Type: model.TokKeyword, Value: buff.String()}, nil
}

// ReadLiteralString reads a literal string (enclosed in parentheses).
func (l *Lexer) ReadLiteralString() (model.Token, error) {

	var buff bytes.Buffer

	depth := 1

	for {

		b, err := l.ReadByte()

		if err != nil {
			return model.Token{}, err
		}

		if b == model.OpenParen {
			depth++
		} else if b == model.CloseParen {
			depth--
			if depth == 0 {
				break
			}
		}

		buff.WriteByte(b)

	}

	return model.Token{Type: model.TokString, Value: buff.String()}, nil

}

// ReadHexaString reads a hexadecimal string (enclosed in angle brackets).
func (l *Lexer) ReadHexaString() (model.Token, error) {

	var buff bytes.Buffer

	for {

		b, err := l.ReadByte()

		if err != nil {
			return model.Token{}, err
		}

		if b == model.GreaterThan {
			break
		}

		if !IsWhiteSpace(b) {
			buff.WriteByte(b)
		}

	}

	return model.Token{Type: model.TokHexString, Value: buff.String()}, nil

}

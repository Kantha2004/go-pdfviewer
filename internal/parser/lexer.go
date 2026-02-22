package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
	"github.com/Kantha2004/go-pdfviewer/internal/util"
)

// Lexer parses a PDF input stream into tokens.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer creates a new Lexer reading from the provided io.Reader.
func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(rd)}
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

		if !util.IsWhiteSpace(b) {
			if err := l.UnReadByte(); err != nil {
				return err
			}
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

		if err := l.UnReadByte(); err != nil {
			return model.Token{}, err
		}
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
		if util.IsNumberChar(b) {
			if err := l.UnReadByte(); err != nil {
				return model.Token{}, err
			}
			return l.ReadNumber()
		}

		if util.IsDelimiter(b) {
			return model.Token{}, fmt.Errorf("unexpected delimiter: %c", b)
		}

		if err := l.UnReadByte(); err != nil {
			return model.Token{}, err
		}
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

		if !util.IsNumberChar(b) {
			if err := l.UnReadByte(); err != nil {
				return model.Token{}, err
			}
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

		if err != nil || util.IsDelimiter(b) || util.IsWhiteSpace(b) {
			if err == nil {
				if err := l.UnReadByte(); err != nil {
					return model.Token{}, err
				}
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

		if err != nil || util.IsDelimiter(b) || util.IsWhiteSpace(b) {
			if err == nil {
				if err := l.UnReadByte(); err != nil {
					return model.Token{}, err
				}
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

		if !util.IsWhiteSpace(b) {
			buff.WriteByte(b)
		}

	}

	return model.Token{Type: model.TokHexString, Value: buff.String()}, nil

}

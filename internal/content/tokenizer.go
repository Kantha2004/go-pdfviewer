package content

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
	"github.com/Kantha2004/go-pdfviewer/internal/util"
)

type CSTokenType int

const (
	CSTokenOperand CSTokenType = iota
	CSTokenOperator
	CSTokenEOF
)

type CSToken struct {
	Type     CSTokenType
	Operand  model.PDFValue
	Operator string
}

func NewContentTokenizer(s io.Reader) *CSTokernizer {
	return &CSTokernizer{
		reader: bufio.NewReader(s),
	}
}

type CSTokernizer struct {
	reader *bufio.Reader
}

func (t *CSTokernizer) peekByte() (byte, error) {
	b, err := t.reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if err := t.reader.UnreadByte(); err != nil {
		return 0, err
	}

	return b, nil
}

func (t *CSTokernizer) skipWhiteSpaceAndComments() error {

	for {
		b, err := t.reader.ReadByte()

		if err != nil {
			return err
		}

		if b == '%' {
			for {
				c, err := t.reader.ReadByte()
				if err != nil || c == '\n' || c == '\r' {
					break
				}
			}
			continue
		}

		if !util.IsWhiteSpace(b) {
			if err := t.reader.UnreadByte(); err != nil {
				return err
			}
			return nil
		}

	}

}

func (t *CSTokernizer) readNumber() (model.PDFNumber, error) {

	var buff bytes.Buffer

	for {
		b, err := t.reader.ReadByte()

		if err != nil {
			return model.PDFNumber(0), err
		}

		if !util.IsNumberChar(b) {
			if err := t.reader.UnreadByte(); err != nil {
				return model.PDFNumber(0), err
			}
			break
		}

		buff.WriteByte(b)

	}

	f, err := strconv.ParseFloat(buff.String(), 64)

	if err != nil {
		return model.PDFNumber(0), err
	}

	return model.PDFNumber(f), nil
}

func (t *CSTokernizer) readName() (model.PDFName, error) {

	var buff bytes.Buffer

	for {
		b, err := t.reader.ReadByte()

		if err != nil || util.IsDelimiter(b) || util.IsWhiteSpace(b) {
			if err == nil {
				if err := t.reader.UnreadByte(); err != nil {
					return "", err
				}
			}
			break
		}

		buff.WriteByte(b)
	}

	return model.PDFName(buff.String()), nil

}

func (t *CSTokernizer) readLiteralString() (model.PDFString, error) {
	var buff bytes.Buffer
	depth := 1

	for {
		b, err := t.reader.ReadByte()
		if err != nil {
			return "", fmt.Errorf("unterminated literal string")
		}

		// Escape handling
		if b == '\\' {
			next, err := t.reader.ReadByte()
			if err != nil {
				return "", fmt.Errorf("unterminated escape sequence")
			}

			switch next {

			// ---- Standard escapes ----
			case 'n':
				buff.WriteByte('\n')
			case 'r':
				buff.WriteByte('\r')
			case 't':
				buff.WriteByte('\t')
			case 'b':
				buff.WriteByte('\b')
			case 'f':
				buff.WriteByte('\f')

			// ---- Escaped parentheses & backslash ----
			case '(', ')', '\\':
				buff.WriteByte(next)

			// ---- Line continuation ----
			case '\r':
				// Consume optional '\n'
				peek, err := t.peekByte()
				if err == nil && peek == '\n' {
					t.reader.ReadByte()
				}
				// Do not write anything (line continuation removed)

			case '\n':
				// Line continuation — ignore entirely

			// ---- Octal escape ----
			case '0', '1', '2', '3', '4', '5', '6', '7':
				octalDigits := []byte{next}

				// Read up to two more octal digits
				for i := 0; i < 2; i++ {
					peek, err := t.peekByte()
					if err != nil {
						break
					}
					if peek >= '0' && peek <= '7' {
						t.reader.ReadByte()
						octalDigits = append(octalDigits, peek)
					} else {
						break
					}
				}

				val, err := strconv.ParseInt(string(octalDigits), 8, 8)
				if err != nil {
					return "", fmt.Errorf("invalid octal escape")
				}

				buff.WriteByte(byte(val))

			// ---- Unknown escape: treat as literal ----
			default:
				buff.WriteByte(next)
			}

			continue
		}

		// ---- Depth tracking ----
		if b == '(' {
			depth++
			buff.WriteByte(b)
			continue
		}

		if b == ')' {
			depth--
			if depth == 0 {
				break
			}
			buff.WriteByte(b)
			continue
		}

		// Normal byte
		buff.WriteByte(b)
	}

	return model.PDFString(buff.String()), nil
}

func (t *CSTokernizer) readHexaString() (model.PDFHexString, error) {
	var hexDigits []byte

	for {
		b, err := t.reader.ReadByte()
		if err != nil {
			return "", fmt.Errorf("unterminated hex string")
		}

		if b == model.GreaterThan {
			break
		}

		if util.IsWhiteSpace(b) {
			continue
		}

		// Validate hex digit
		if !util.IsHexDigit(b) {
			return "", fmt.Errorf("invalid hex digit: %c", b)
		}

		hexDigits = append(hexDigits, b)
	}

	// If odd number of digits, pad with '0'
	if len(hexDigits)%2 != 0 {
		hexDigits = append(hexDigits, '0')
	}

	// Convert hex pairs to bytes
	var decoded bytes.Buffer

	for i := 0; i < len(hexDigits); i += 2 {
		pair := string(hexDigits[i : i+2])

		val, err := strconv.ParseUint(pair, 16, 8)
		if err != nil {
			return "", fmt.Errorf("invalid hex pair: %s", pair)
		}

		decoded.WriteByte(byte(val))
	}

	return model.PDFHexString(decoded.String()), nil
}

func (t *CSTokernizer) readValue() (model.PDFValue, error) {
	b, err := t.reader.ReadByte()

	if err != nil {
		return nil, err
	}

	switch b {

	case model.Solidus:
		return t.readName()

	case model.OpenParen:
		return t.readLiteralString()

	case model.LessThan:
		return t.readHexaString()

	case model.OpenLBracket:
		return t.readArray()

	default:
		if util.IsNumberChar(b) {
			if err := t.reader.UnreadByte(); err != nil {
				return nil, err
			}
			return t.readNumber()
		}

		if util.IsDelimiter(b) {
			return nil, fmt.Errorf("unexpected delimiter: %c", b)
		}

		if err := t.reader.UnreadByte(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("not a valid operand")
	}
}

func (t *CSTokernizer) readArray() (model.PDFArray, error) {
	var arr []model.PDFValue

	for {

		if err := t.skipWhiteSpaceAndComments(); err != nil {
			return arr, err
		}

		b, err := t.reader.ReadByte()

		if err != nil {
			return arr, err
		}

		if b == model.CloseLBracket {
			break
		}

		err = t.reader.UnreadByte()

		if err != nil {
			return arr, err
		}

		operand, err := t.readValue()

		if err != nil {
			return arr, err
		}

		arr = append(arr, operand)

	}

	return model.PDFArray(arr), nil

}

func (t *CSTokernizer) handleOperatorEOF(err error, buff *bytes.Buffer) (CSToken, error) {
	if err != io.EOF {
		return CSToken{}, err
	}

	if buff.Len() == 0 {
		return CSToken{Type: CSTokenEOF}, nil
	}

	return CSToken{
		Type:     CSTokenOperator,
		Operator: buff.String(),
	}, nil

}

func (t *CSTokernizer) readOperator() (CSToken, error) {
	var buff bytes.Buffer

	for {
		b, err := t.reader.ReadByte()
		if err != nil {
			return t.handleOperatorEOF(err, &buff)
		}

		if util.IsWhiteSpace(b) || util.IsDelimiter(b) {
			if err := t.reader.UnreadByte(); err != nil {
				return CSToken{}, err
			}
			break
		}

		buff.WriteByte(b)
	}

	if buff.Len() == 0 {
		return CSToken{}, fmt.Errorf("empty operator token")
	}

	return CSToken{
		Type:     CSTokenOperator,
		Operator: buff.String(),
	}, nil
}

func (t *CSTokernizer) operand(operandValue model.PDFValue) CSToken {
	return CSToken{Type: CSTokenOperand, Operand: operandValue}
}

func (t *CSTokernizer) NextToken() (CSToken, error) {
	if err := t.skipWhiteSpaceAndComments(); err != nil {
		if err == io.EOF {
			return CSToken{Type: CSTokenEOF}, nil
		}
		return CSToken{}, err
	}

	b, err := t.peekByte()

	if err != nil {
		if err == io.EOF {
			return CSToken{Type: CSTokenEOF}, nil
		}
		return CSToken{}, err
	}

	if util.IsValueStart(b) {
		operand, err := t.readValue()
		if err != nil {
			return CSToken{}, err
		}
		return t.operand(operand), nil
	} else {
		return t.readOperator()
	}

}

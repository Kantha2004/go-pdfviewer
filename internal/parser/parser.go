package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

type Parser struct {
	l   *Lexer
	buf *model.Token
}

func NewParser(l *Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) next() (model.Token, error) {
	if p.buf != nil {
		t := *p.buf
		p.buf = nil
		return t, nil
	}

	return p.l.NextToken()
}

func (p *Parser) unread(t model.Token) {
	p.buf = &t
}

func (p *Parser) Parse() (model.PDFValue, error) {

	tok, err := p.next()

	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case model.TokNumber:
		return p.parseNumberOrRef(tok)
	case model.TokName:
		return model.PDFName(tok.Value), nil
	case model.TokString:
		return model.PDFString(tok.Value), nil
	case model.TokHexString:
		return model.PDFHexString(tok.Value), nil
	case model.TokKeyword:
		return p.parseKeyword(tok)
	case model.TokArrayStart:
		return p.parseArray()
	case model.TokDictStart:
		return p.parseDict()
	default:
		return nil, fmt.Errorf("unexpected token: type=%v value=%q", tok.Type, tok.Value)
	}
}

func (p *Parser) parseNumberOrRef(first model.Token) (model.PDFValue, error) {
	n1, err := strconv.Atoi(first.Value)

	if err != nil {
		// If not an integer, try parsing as float
		f, err := strconv.ParseFloat(first.Value, 64)
		if err != nil {
			return nil, err
		}
		return model.PDFNumber(f), nil
	}

	// It's an integer, could be start of Indirect Ref: Int Int R
	tok2, err := p.next()
	if err != nil {
		// End of stream or error, return the number
		return model.PDFNumber(float64(n1)), nil
	}

	if tok2.Type != model.TokNumber {
		p.unread(tok2)
		return model.PDFNumber(float64(n1)), nil
	}

	n2, err := strconv.Atoi(tok2.Value)
	if err != nil {
		// Second token is a number but not an integer (e.g. 1 2.5), so not a ref.
		// Unread tok2 and return n1.
		p.unread(tok2)
		return model.PDFNumber(float64(n1)), nil
	}

	tok3, err := p.next()
	if err != nil {
		p.unread(tok2)
		return model.PDFNumber(float64(n1)), nil
	}

	if tok3.Type == model.TokKeyword && tok3.Value == "R" {
		return model.PDFIndirectRef{
			ObjectNumber: n1,
			Generation:   n2,
		}, nil
	}

	p.unread(tok3)

	return model.PDFNumber(float64(n1)), nil
}

func (p *Parser) parseArray() (model.PDFValue, error) {

	arr := make([]model.PDFValue, 0, 4)

	for {

		tok, err := p.next()

		if err != nil {
			return nil, err
		}

		if tok.Type == model.TokArrayEnd {
			break
		}

		if tok.Type == model.TokEOF {
			return nil, fmt.Errorf("unterminated array")
		}

		p.unread(tok)

		val, err := p.Parse()

		if err != nil {
			return nil, err
		}

		arr = append(arr, val)

	}

	return arr, nil

}

func (p *Parser) parseDict() (model.PDFValue, error) {

	dict := make(map[string]model.PDFValue)

	for {
		tok, err := p.next()

		if err != nil {
			return nil, err
		}

		if tok.Type == model.TokDictEnd {
			break
		}

		if tok.Type != model.TokName {
			return nil, fmt.Errorf("dictionary key must be name, got %v", tok)
		}

		if tok.Type == model.TokEOF {
			return nil, fmt.Errorf("unterminated dictionary")
		}

		val, err := p.Parse()

		if err != nil {
			return nil, err
		}

		dict[tok.Value] = val

	}

	return dict, nil

}

func (p *Parser) parseKeyword(tok model.Token) (model.PDFValue, error) {

	switch tok.Value {
	case "true":
		return model.PDFBoolean(true), nil
	case "false":
		return model.PDFBoolean(false), nil
	case "null":
		return model.PDFNull{}, nil
	default:
		return nil, fmt.Errorf("unexpected keyword: %s", tok.Value)
	}

}

func (p *Parser) ParseObject() (*model.PDFObject, error) {
	// ---- object number ----
	tok, err := p.next()
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, err
	}

	if tok.Type == model.TokEOF {
		return nil, io.EOF
	}

	if tok.Type != model.TokNumber {
		return nil, fmt.Errorf("expected object number, got %v", tok)
	}

	objNum, err := strconv.Atoi(tok.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid object number %q", tok.Value)
	}

	// ---- generation number ----
	tok, err = p.next()
	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokNumber {
		return nil, fmt.Errorf("expected generation number, got %v", tok)
	}

	genNum, err := strconv.Atoi(tok.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid generation number %q", tok.Value)
	}

	// ---- expect 'obj' ----
	tok, err = p.next()
	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokKeyword || tok.Value != model.ObjectStart {
		return nil, fmt.Errorf("expected 'obj', got %v", tok)
	}

	// ---- parse object value ----
	val, err := p.Parse()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("unterminated object %d %d", objNum, genNum)
		}
		return nil, err
	}

	// ---- expect 'endobj' ----
	tok, err = p.next()
	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokKeyword || tok.Value != model.ObjectEnd {
		return nil, fmt.Errorf("expected 'endobj', got %v", tok)
	}

	return &model.PDFObject{
		Number: objNum,
		Gen:    genNum,
		Value:  val,
	}, nil
}

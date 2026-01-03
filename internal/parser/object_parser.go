package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

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

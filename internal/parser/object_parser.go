package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func (p *Parser) ResolveStreamLength(dict model.PDFDict) (int, error) {
	// ---- direct length ----
	if n, ok := dict["Length"].(model.PDFNumber); ok {
		return int(n), nil
	}

	// ---- indirect length ----
	ref, ok := dict["Length"].(model.PDFIndirectRef)
	if !ok {
		return 0, fmt.Errorf("stream missing valid /Length")
	}

	obj, ok := p.objects.Get(ref.ObjectNumber, ref.Generation)
	if !ok {
		return 0, fmt.Errorf("stream /Length reference %d %d not found",
			ref.ObjectNumber, ref.Generation)
	}

	num, ok := obj.Value.(model.PDFNumber)
	if !ok {
		return 0, fmt.Errorf("stream /Length object %d %d is not a number",
			ref.ObjectNumber, ref.Generation)
	}

	return int(num), nil
}

func (p *Parser) ConsumeEOL() error {
	newLine, err := p.l.ReadByte()

	if err != nil {
		return err
	}

	if newLine == '\r' {
		newLine, err := p.l.ReadByte()

		if err != nil {
			return err
		}

		if newLine != '\n' {
			return fmt.Errorf("expected a newline")
		}
	} else if newLine != '\n' {
		return fmt.Errorf("expected a newline")
	}

	return nil
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

	// ---- expect 'endobj' or 'stream' ----
	tok, err = p.next()
	if err != nil {
		return nil, err
	}

	// this is 'stream' the prebious value should be a dict
	// ---- stream ----
	if tok.Type == model.TokKeyword && tok.Value == model.StreamStart {

		valDict, ok := val.(model.PDFDict)
		if !ok {
			return nil, fmt.Errorf("expected Dict value before stream")
		}

		length, err := p.ResolveStreamLength(valDict)
		if err != nil {
			return nil, err
		}

		// Mandatory EOL after 'stream'
		if err := p.ConsumeEOL(); err != nil {
			return nil, err
		}

		data := make([]byte, length)
		for i := 0; i < length; i++ {
			b, err := p.l.ReadByte()
			if err != nil {
				return nil, err
			}
			data[i] = b
		}

		tok, err = p.next()
		if err != nil {
			return nil, err
		}

		if tok.Type != model.TokKeyword || tok.Value != model.StreamEnd {
			return nil, fmt.Errorf("expected 'endstream', got %v", tok)
		}

		val = model.PDFStream{
			Dict: valDict,
			Data: data,
		}

		tok, err = p.next()
		if err != nil {
			return nil, err
		}
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

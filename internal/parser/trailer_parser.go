package parser

import (
	"fmt"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func (p *Parser) ParseTrailer() (model.PDFDict, error) {

	tok, err := p.next()

	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokKeyword || tok.Value != model.Trailer {
		return nil, fmt.Errorf("not valid trailer")
	}

	v, err := p.Parse()
	if err != nil {
		return nil, err
	}

	dict, ok := v.(model.PDFDict)
	if !ok {
		return nil, fmt.Errorf("trailer is not a dictionary")
	}

	return dict, nil

}

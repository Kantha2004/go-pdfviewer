package parser

import (
	"fmt"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func (p *Parser) ParseTrailer() (model.PDFValue, error) {

	tok, err := p.next()

	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokKeyword || tok.Value != model.Trailer {
		return nil, fmt.Errorf("not valid trailer")
	}

	return p.Parse()

}

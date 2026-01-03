package parser

import (
	"fmt"
	"io"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

type Document struct {
	Objects *ObjectTable
	XRef    *model.XRefTable
	Trailer model.PDFValue
}

func (p *Parser) ParseDocument() (*Document, error) {

	doc := &Document{
		Objects: NewObjectTable(),
	}

	for {

		tok, err := p.next()

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if tok.Type == model.TokKeyword && tok.Value == model.XRef {
			p.unread(tok)
			xref, err := p.ParseXRef()
			if err != nil {
				return nil, err
			}

			doc.XRef = xref

			trailer, err := p.ParseTrailer()

			if err != nil {
				return nil, err
			}

			doc.Trailer = trailer

			fmt.Printf("ParsedXref %+v\n", doc.XRef)
			fmt.Printf("ParsedTrailer %+v\n", doc.Trailer)
			break
		}

		p.unread(tok)

		obj, err := p.ParseObject()

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		doc.Objects.Add(obj)

	}

	return doc, nil
}

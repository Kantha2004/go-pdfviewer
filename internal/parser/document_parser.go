package parser

import (
	"fmt"
	"io"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

type Document struct {
	Objects *ObjectTable
	XRef    *model.XRefTable
	Trailer model.PDFDict
	Catalog *model.PDFObject
	Pages   []*model.PDFObject
}

func (doc *Document) ResolveCatalog() error {

	rootVal, ok := doc.Trailer["Root"]
	if !ok {
		return fmt.Errorf("missing /Root in trailer")
	}

	rootRef, ok := rootVal.(model.PDFIndirectRef)
	if !ok {
		return fmt.Errorf("/Root is not an indirect reference")
	}

	catalogObj, ok := doc.Objects.Get(rootRef.ObjectNumber, rootRef.Generation)
	if !ok {
		return fmt.Errorf("root catalog object %d %d not found", rootRef.ObjectNumber, rootRef.Generation)
	}

	doc.Catalog = catalogObj
	return nil

}

func (doc *Document) ResolveEachPage(objNum int, gen int) error {

	pages, ok := doc.Objects.Get(objNum, gen)

	if !ok {
		return fmt.Errorf("unable to find page object")
	}

	pagesDict, ok := pages.Value.(model.PDFDict)

	if !ok {
		return fmt.Errorf("pages is not a dict")
	}

	pagesType, ok := pagesDict["Type"].(model.PDFName)
	if !ok {
		return fmt.Errorf("page node missing /Type")
	}

	switch pagesType {

	case model.PagesType:
		kids, ok := pagesDict["Kids"].(model.PDFArray)

		if !ok {
			return fmt.Errorf("kids is not an array %v", pagesDict["Kids"])
		}

		for _, p := range kids {

			pRef, ok := p.(model.PDFIndirectRef)

			if !ok {
				return fmt.Errorf("not a valid ref: %v", p)
			}

			err := doc.ResolveEachPage(pRef.ObjectNumber, pRef.Generation)

			if err != nil {
				return err
			}

		}

	case model.PageType:
		doc.Pages = append(doc.Pages, pages)

	default:
		return fmt.Errorf("invalid type found instead of %v or %v found %v", model.PagesType, model.PageType, pagesDict["Type"])

	}

	return nil
}

func (doc *Document) ResolvePages() error {
	doc.Pages = make([]*model.PDFObject, 0)

	catalogDict, ok := doc.Catalog.Value.(model.PDFDict)

	if !ok {
		return fmt.Errorf("catalog is not a Dict")
	}

	pages, ok := catalogDict["Pages"].(model.PDFIndirectRef)

	if !ok {
		return fmt.Errorf("pages is not in catalog %v", catalogDict)
	}

	err := doc.ResolveEachPage(pages.ObjectNumber, pages.Generation)

	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) ParseDocument() (*Document, error) {

	doc := &Document{
		Objects: NewObjectTable(),
	}

	p.objects = doc.Objects

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

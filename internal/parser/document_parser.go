package parser

import (
	"fmt"
	"io"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

type Page struct {
	// Identity
	ObjectNumber int
	Generation   int

	Parent model.PDFIndirectRef
	Dict   model.PDFDict

	// Structural properties (resolved)
	MediaBox  [4]float64
	Resources model.PDFDict

	// Raw content streams (resolved but not interpreted)
	Streams []model.PDFStream

	// Parsed instructions (filled later)
	// Instructions []Instruction
}

type Document struct {
	Objects       *ObjectTable
	XRef          *model.XRefTable
	Trailer       model.PDFDict
	Catalog       *model.PDFObject
	Pages         []*model.PDFObject
	ResolvedPages []Page
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

func (doc *Document) ResolvePageContents() error {

	doc.ResolvedPages = make([]Page, 0)

	for _, page := range doc.Pages {
		pageDict, ok := page.Value.(model.PDFDict)

		if ok != true {
			return fmt.Errorf("page value is not a dict %v", page.Value)
		}

		// Contents
		contents := pageDict["Contents"]

		if contents == nil {
			return fmt.Errorf("contents is missing from page %v", pageDict)
		}

		pageObj := Page{
			Dict:         pageDict,
			ObjectNumber: page.Number,
			Generation:   page.Gen,
			Streams:      make([]model.PDFStream, 0),
		}

		parentRef, ok := pageDict["Parent"].(model.PDFIndirectRef)
		if ok {
			pageObj.Parent = parentRef
		}

		switch contents := contents.(type) {

		case model.PDFIndirectRef:

			err := doc.addStreamToPage(contents, &pageObj)
			if err != nil {
				return err
			}

		case model.PDFArray:
			for _, content := range contents {

				ref, ok := content.(model.PDFIndirectRef)

				if !ok {
					return fmt.Errorf("page content is not the required type %t", content)
				}

				err := doc.addStreamToPage(ref, &pageObj)
				if err != nil {
					return err
				}

			}

		default:
			return fmt.Errorf("page content is not the required type %t", pageDict["Contents"])

		}

		// MediaBox
		if err := doc.addMediaBoxToPage(&pageObj); err != nil {
			return err
		}

		// Resources
		if err := doc.addResourcesBoxToPage(&pageObj); err != nil {
			return err
		}

		doc.ResolvedPages = append(doc.ResolvedPages, pageObj)

	}

	return nil
}

func (doc *Document) addMediaBoxToPage(page *Page) error {

	mediabox, err := doc.resolveInheritance(page, "MediaBox")

	if err != nil {
		return err
	}

	mediaBox, ok := mediabox.(model.PDFArray)

	if !ok {
		return fmt.Errorf("media Box should of type array but got %t", mediaBox)
	}

	if mediaBox != nil {
		if len(mediaBox) != 4 {
			return fmt.Errorf("mediaBox size should be 4 but len is %d value :%v", len(mediaBox), mediaBox)
		}

		var box [4]float64

		for i, number := range mediaBox {

			num, ok := number.(model.PDFNumber)
			if !ok {
				return fmt.Errorf("mediaBox values should of type number but %t", number)
			}

			box[i] = float64(num)
		}

		page.MediaBox = box
	}

	return nil

}

func (doc *Document) addResourcesBoxToPage(page *Page) error {

	resource, err := doc.resolveInheritance(page, "Resources")

	if err != nil {
		return err
	}

	resources, ok := resource.(model.PDFDict)
	if !ok {
		return fmt.Errorf("resources should be dictionary but got %T", resource)
	}

	page.Resources = resources

	return nil

}

func (doc *Document) addStreamToPage(content model.PDFIndirectRef, page *Page) error {
	object, ok := doc.Objects.GetObjectValue(content.ObjectNumber, content.Generation)

	if !ok {
		return fmt.Errorf("page content is not found")
	}

	stream, ok := object.(model.PDFStream)

	if !ok {
		return fmt.Errorf("page content is not stream %v", object)
	}

	page.Streams = append(page.Streams, stream)

	return nil
}

func (doc *Document) resolveInheritance(page *Page, key string) (any, error) {

	currentDict := page.Dict
	currentParent := page.Parent

	for {

		if value, ok := currentDict[key]; ok {
			return value, nil
		}

		if currentParent.ObjectNumber == 0 {
			break
		}

		parentObj, ok := doc.Objects.GetObjectValue(currentParent.ObjectNumber, currentParent.Generation)
		if !ok {
			return nil, fmt.Errorf("parent object %d %d not found",
				currentParent.ObjectNumber, currentParent.Generation)
		}

		parentDict, ok := parentObj.(model.PDFDict)

		if !ok {
			return nil, fmt.Errorf("the parent object %d %d is not a dictionary", currentParent.ObjectNumber, currentParent.Generation)
		}

		currentDict = parentDict

		if parentRef, ok := parentDict["Parent"].(model.PDFIndirectRef); ok {
			currentParent = parentRef
		} else {
			currentParent = model.PDFIndirectRef{}
		}

	}

	return nil, fmt.Errorf("property %s not found in inheritance chain", key)
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

package parser

import (
	"fmt"
	"strconv"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

func (p *Parser) ParseXRef() (*model.XRefTable, error) {

	xRefTable := make(model.XRefTable)
	tok, err := p.next()

	if err != nil {
		return nil, err
	}

	if tok.Type != model.TokKeyword || tok.Value != model.XRef {
		return nil, fmt.Errorf("Not a valid xref")
	}

	for {

		objTok, err := p.next()

		if err != nil {
			return nil, err
		}

		if objTok.Type == model.TokKeyword && objTok.Value == model.Trailer {
			p.unread(objTok)
			break
		}

		if objTok.Type != model.TokNumber {
			return nil, fmt.Errorf("not a valid object")
		}

		objectIndex, err := strconv.Atoi(objTok.Value)

		if err != nil {
			return nil, err
		}

		xRefCountTok, err := p.next()

		if err != nil {
			return nil, err
		}

		if xRefCountTok.Type != model.TokNumber {
			return nil, fmt.Errorf("not a valid object")
		}

		xRefCount, err := strconv.Atoi(xRefCountTok.Value)

		if err != nil {
			return nil, err
		}

		for i := range xRefCount {

			offsetTok, err := p.next()

			if err != nil {
				return nil, err
			}

			if offsetTok.Type != model.TokNumber {
				return nil, fmt.Errorf("not a valid offset")
			}

			offset, err := strconv.Atoi(offsetTok.Value)

			if err != nil {
				return nil, err
			}

			//

			genTok, err := p.next()

			if err != nil {
				return nil, err
			}

			if genTok.Type != model.TokNumber {
				return nil, fmt.Errorf("not a valid gen")
			}

			gen, err := strconv.Atoi(genTok.Value)

			if err != nil {
				return nil, err
			}

			//

			isUsingTok, err := p.next()

			if err != nil {
				return nil, err
			}

			if isUsingTok.Type != model.TokKeyword || (isUsingTok.Value != model.ObjectInUse && isUsingTok.Value != model.ObjectFree) {
				return nil, fmt.Errorf("not a valid state")
			}

			objNum := objectIndex + i

			xRefTable[objNum] = model.XRefEntry{
				Offset:     offset,
				Generation: gen,
				InUse:      isUsingTok.Value == model.ObjectInUse,
			}

		}

	}

	return &xRefTable, nil
}

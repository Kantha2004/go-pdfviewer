package parser

import (
	"github.com/Kantha2004/go-pdfviewer/internal/model"
)

type ObjectTable struct {
	Ref map[int]map[int]*model.PDFObject
}

func NewObjectTable() *ObjectTable {
	return &ObjectTable{
		Ref: make(map[int]map[int]*model.PDFObject),
	}
}

func (o *ObjectTable) Add(m *model.PDFObject) {
	if o.Ref[m.Number] == nil {
		o.Ref[m.Number] = make(map[int]*model.PDFObject)
	}
	o.Ref[m.Number][m.Gen] = m
}

func (o *ObjectTable) Get(objectNum int, gen int) (*model.PDFObject, bool) {
	if objs := o.Ref[objectNum]; objs != nil {
		if obj := objs[gen]; obj != nil {
			return obj, true
		}
	}

	return nil, false
}

package model

type PDFValue any

type PDFNull struct{}

type PDFBoolean bool

type PDFNumber float64

type PDFName string

type PDFString string

type PDFHexString string

type PDFArray []PDFArray

type PDFDict map[string]PDFValue

type PDFIndirectRef struct {
	ObjectNumber int
	Generation   int
}

type PDFObject struct {
	Number int
	Gen    int
	Value  PDFValue
}

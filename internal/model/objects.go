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

type XRefEntry struct {
	Offset     int
	Generation int
	InUse      bool
}

type XRefTable map[int]XRefEntry

const (
	ObjectInUse string = "n"
	ObjectFree  string = "f"
)

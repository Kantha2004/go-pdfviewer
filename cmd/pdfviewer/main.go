package main

import (
	"fmt"
	"os"

	"github.com/Kantha2004/go-pdfviewer/internal/parser"
)

func main() {
	pdfFile, err := os.Open("testdata/minimal.pdf")

	if err != nil {
		fmt.Printf("Error opening pdf file %q", err)
		return
	}
	defer pdfFile.Close()

	l := parser.NewLexer(pdfFile)

	p := parser.NewParser(l)

	doc, err := p.ParseDocument()

	if err != nil {
		fmt.Printf("Error while parsing pdf file %v\n", err)
		return
	}

	err = doc.ResolveCatalog()

	if err != nil {
		fmt.Printf("Error while resolving catalog: %v\n", err)
		return
	}

	err = doc.ResolvePages()

	if err != nil {
		fmt.Printf("Error while resolving pages: %v\n", err)
		return
	}

	fmt.Println("PDF file parsed successfully!")

	for i, obj := range doc.Pages {
		fmt.Printf("Page %d: %#v\n", i, obj)
	}

	// for objNum, gens := range doc.Objects.Ref {
	// 	for gen, obj := range gens {
	// 		fmt.Printf("Object %d %d -> %T\n", objNum, gen, obj.Value)
	// 	}
	// }
}

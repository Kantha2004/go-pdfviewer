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

	// objTable := parser.NewObjectTable()

	_, err = p.ParseDocument()

	if err != nil {
		fmt.Printf("Eroor while parsing pdf file %v", err)
		return
	}

	// fmt.

	// for {
	// 	val, err := p.ParseObject()

	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		fmt.Printf("Error: %v\n", err)
	// 		return
	// 	}

	// 	objTable.Add(val)
	// 	fmt.Printf("Parsed Value: %+v\n", val)
	// 	ob, _ := objTable.Get(val.Number, val.Gen)
	// 	fmt.Printf("Parsed In Object table: %+v\n", ob)
	// }

	// for objNum, gens := range objTable.Ref {
	// 	for gen, obj := range gens {
	// 		fmt.Printf("Object %d %d -> %T\n", objNum, gen, obj.Value)
	// 	}
	// }
}

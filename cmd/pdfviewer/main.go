package main

import (
	"fmt"
	"io"
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

	for {
		val, err := p.Parse()

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Parsed Value: %+v\n", val)

	}

}

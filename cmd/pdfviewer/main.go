package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Kantha2004/go-pdfviewer/internal/model"
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

	for {
		tok, err := l.NextToken()

		if err != nil {
			if err == io.EOF || tok.Type == model.TokEOF {
				break
			}
			fmt.Printf("Error: %v\n", err)
			return
		}

		if tok.Type == model.TokEOF {
			break
		}

		fmt.Printf("Token: %v, Value: %q\n", tok.Type, tok.Value)

	}

}

# go-pdfviewer

A from-scratch PDF parser and renderer written in Go.

This project is an educational and engineering-focused implementation of a PDF
engine that:
- Parses PDF files at the byte level
- Implements its own lexer, parser, and object model
- Interprets PDF content streams
- Renders vector graphics and text without relying on existing PDF engines

The goal is correctness, clarity, and deep understanding — not shortcuts.

---

## Project Goals

- Implement a PDF lexer (tokenizer) from first principles
- Build a recursive PDF value parser
- Resolve indirect objects and cross-reference tables
- Traverse the PDF page tree
- Interpret PDF graphics operators
- Rasterize vector graphics into pixel images
- Gradually add text, fonts, and images

This project intentionally avoids using existing PDF rendering libraries.

---

## Project Structure

```text
go-pdfviewer/
├── cmd/pdfviewer/          # Entry point (CLI)
├── internal/
│   ├── parser/             # Lexer, parser, xref, stream decoding
│   ├── model/              # PDF object model
│   ├── graphics/           # Graphics state & content interpreter
│   ├── render/             # Rasterizer & image output
│   └── util/               # Helpers (math, IO)
├── testdata/               # Minimal and test PDFs
├── go.mod
└── README.md
```

---

## Current Status

- [x] Project structure
- [ ] PDF lexer
- [ ] PDF value parser
- [ ] Indirect object parsing
- [ ] XRef resolution
- [ ] Page tree traversal
- [ ] Rendering engine

---

## How to Run (Lexer Debug)

Temporary debug example in `cmd/pdfviewer/main.go`:

```go
for {
    tok, err := lexer.NextToken()
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", tok)
    if tok.Type == TokEOF {
        break
    }
}
```

Run:
```bash
go run ./cmd/pdfviewer
```

---

## Design Principles

- Byte-accurate parsing (no line-based assumptions)
- Strict separation of concerns
- Forgiving parsing for malformed PDFs
- Spec-aligned behavior (ISO 32000-1)

---

## What This Is Not

- Not a wrapper around MuPDF or PDFium
- Not a production-ready PDF viewer
- Not a shortcut-based implementation

This is a learning-first, engine-level project.

---

## Roadmap (High-Level)

1. Lexer
2. Value parser
3. Object & XRef parsing
4. Page tree resolution
5. Content stream interpreter
6. Rasterizer
7. Text and fonts
8. Images and advanced features

---

## License

MIT (or choose later)

---

## Author

Kanthakumar
GitHub: https://github.com/Kantha2004

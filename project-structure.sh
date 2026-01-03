mkdir -p go-pdfviewer/{cmd/pdfviewer,internal/{parser,model,graphics,render,util},testdata} \
&& touch \
go-pdfviewer/cmd/go-pdfviewer/main.go \
go-pdfviewer/internal/parser/{lexer.go,tokens.go,parser.go,xref.go,stream.go} \
go-pdfviewer/internal/model/{objects.go,document.go,resources.go} \
go-pdfviewer/internal/graphics/{state.go,operators.go,interpreter.go,text.go} \
go-pdfviewer/internal/render/{rasterizer.go,path.go,image.go,transform.go} \
go-pdfviewer/internal/util/{reader.go,math.go} \
go-pdfviewer/testdata/minimal.pdf \
go-pdfviewer/{go.mod,README.md}

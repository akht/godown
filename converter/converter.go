package converter

import (
	"bufio"
	"bytes"
	"godown/evaluator"
	"godown/lexer"
	"godown/parser"
	"io"
)

func Convert(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	var buf bytes.Buffer
	for {
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		buf.WriteString(scanner.Text())
		buf.WriteString("\n")
	}

	input := buf.String()
	l := lexer.New(input)
	p := parser.New(l)

	document := p.ParseDocument()

	evaluated := evaluator.Eval(document)
	io.WriteString(out, evaluated.Render())
}

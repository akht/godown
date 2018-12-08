package repl

import (
	"bufio"
	"fmt"
	"godown/evaluator"
	"godown/lexer"
	"godown/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		document := p.ParseDocument()

		evaluated := evaluator.Eval(document)
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }

		io.WriteString(out, evaluated.Inspect())
		// io.WriteString(out, "\n")
	}
}

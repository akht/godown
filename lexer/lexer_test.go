package lexer

import (
	"godown/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `# h1
## h2
###### h6
#text1
######text6
*em*
**strong**
***emstrong***
- list
1. list
999. list
juststring


# *emphasis heading*

~~strikethrough~~
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IGETA, "#"},
		{token.SPACE, " "},
		{token.TEXT, "h1"},
		{token.CR, "\n"},

		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.SPACE, " "},
		{token.TEXT, "h2"},
		{token.CR, "\n"},

		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.SPACE, " "},
		{token.TEXT, "h6"},
		{token.CR, "\n"},

		{token.IGETA, "#"},
		{token.TEXT, "text1"},
		{token.CR, "\n"},

		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.IGETA, "#"},
		{token.TEXT, "text6"},
		{token.CR, "\n"},

		{token.ASTERISK, "*"},
		{token.TEXT, "em"},
		{token.ASTERISK, "*"},
		{token.CR, "\n"},

		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.TEXT, "strong"},
		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.CR, "\n"},

		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.TEXT, "emstrong"},
		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.ASTERISK, "*"},
		{token.CR, "\n"},

		{token.HYPHEN, "-"},
		{token.SPACE, " "},
		{token.TEXT, "list"},
		{token.CR, "\n"},

		{token.TEXT, "1."},
		{token.SPACE, " "},
		{token.TEXT, "list"},
		{token.CR, "\n"},

		{token.TEXT, "999."},
		{token.SPACE, " "},
		{token.TEXT, "list"},
		{token.CR, "\n"},

		{token.TEXT, "juststring"},
		{token.CR, "\n"},

		{token.CR, "\n"},
		{token.CR, "\n"},

		{token.IGETA, "#"},
		{token.SPACE, " "},
		{token.ASTERISK, "*"},
		{token.TEXT, "emphasis"},
		{token.SPACE, " "},
		{token.TEXT, "heading"},
		{token.ASTERISK, "*"},
		{token.CR, "\n"},

		{token.CR, "\n"},

		{token.TILDE, "~"},
		{token.TILDE, "~"},
		{token.TEXT, "strikethrough"},
		{token.TILDE, "~"},
		{token.TILDE, "~"},
		{token.CR, "\n"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] = tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestBackQuoteToken(t *testing.T) {
	input := "`inlinecode`"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BACKQUOTE, "`"},
		{token.TEXT, "inlinecode"},
		{token.BACKQUOTE, "`"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] = tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestCodeBlock(t *testing.T) {
	input := "```go\n"
	input += "func main() {\n"
	input += "    fmt.Printf(\"Hello, world!\")\n"
	input += "}\n"
	input += "```"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BACKQUOTE, "`"},
		{token.BACKQUOTE, "`"},
		{token.BACKQUOTE, "`"},
		{token.TEXT, "go"},
		{token.CR, "\n"},
		{token.TEXT, "func"},
		{token.SPACE, " "},
		{token.TEXT, "main()"},
		{token.SPACE, " "},
		{token.TEXT, "{"},
		{token.CR, "\n"},
		{token.SPACE, " "},
		{token.SPACE, " "},
		{token.SPACE, " "},
		{token.SPACE, " "},
		{token.TEXT, "fmt.Printf(\"Hello,"},
		{token.SPACE, " "},
		{token.TEXT, "world!\")"},
		{token.CR, "\n"},
		{token.TEXT, "}"},
		{token.CR, "\n"},
		{token.BACKQUOTE, "`"},
		{token.BACKQUOTE, "`"},
		{token.BACKQUOTE, "`"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] = tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

}

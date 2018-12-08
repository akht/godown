package evaluator

import (
	"godown/lexer"
	"godown/object"
	"godown/parser"
	"testing"
)

func TestHeadingObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"# heading1", "<h1>heading1</h1>\n"},
		{"## heading2", "<h2>heading2</h2>\n"},
		{"### heading3", "<h3>heading3</h3>\n"},
		{"#### heading4", "<h4>heading4</h4>\n"},
		{"##### heading5", "<h5>heading5</h5>\n"},
		{"###### heading6", "<h6>heading6</h6>\n"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testHeadingObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Document {
	l := lexer.New(input)
	p := parser.New(l)
	document := p.ParseDocument()

	return Eval(document)
}

func testHeadingObject(t *testing.T, doc object.Document, expected string) bool {
	result, ok := doc.Objects[0].(*object.Heading)
	if !ok {
		t.Errorf("object is not Heading. got=%T (%+v)", doc.Objects[0], doc.Objects[0])
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}

	return true
}

func TestDocument(t *testing.T) {
	input := `
# godwon Markdown Parser in Go

## Markdown Spec

- Heading
- Emphasis
- **(em)
- ****(strong)
- ******(em strong)
- Strikethrough
- ~~
- List(DISC)
- List(Decimal)
- Quote
- Horizontal Line
- Code Block
- Inline Code
- Link
- Imange
- Table
`

	expected := `<h1><em>text</em>1MIDASHI1</h1>
<h2>M2</h2>
<p><em>2text</em><em>text2</em></p>
<h2>Heading<em>2</em> <em>text</em></h2>
<p>3text 999 hoge</p>
<p>
<ul>
<li>buy milk</li>
<li>mail tanaka</li>
<li>meeting as 15</li>
</ul>
</p>
<h3>h3<strong>!!!</strong></h3>
<p>This is a text.</p>
`

	evaluated := testEval(input)
	if expected != evaluated.Inspect() {
		t.Errorf("object has wrong value. got=%s, want=%s",
			evaluated.Inspect(), expected)
	}
}

func TestDocument2(t *testing.T) {
	input := "# *text*1MIDASHI1"

	expected := "<h1><em>text</em>1MIDASHI1</h1>\n"

	evaluated := testEval(input)
	if expected != evaluated.Inspect() {
		t.Errorf("object has wrong value. got=%s, want=%s",
			evaluated.Inspect(), expected)
	}
}

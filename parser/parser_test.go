package parser

import (
	"godown/lexer"
	"testing"
)

func TestHeading(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"# text",
			"<h1>text</h1>\n",
		},
		{
			"## text",
			"<h2>text</h2>\n",
		},
		{
			"###### text",
			"<h6>text</h6>\n",
		},
		{
			"# *text*",
			"<h1><em>text</em></h1>\n",
		},
		{
			"# **text**",
			"<h1><strong>text</strong></h1>\n",
		},
		{
			"# ***-text-***",
			"<h1><strong><em>-text-</em></strong></h1>\n",
		},
		{
			"## Heading*2*",
			"<h2>Heading<em>2</em></h2>\n",
		},
		{
			"## -text-",
			"<h2>-text-</h2>\n",
		},
		{
			"## - text *-*",
			"<h2>- text <em>-</em></h2>\n",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		document := p.ParseDocument()

		actual := document.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// ひとまず1行だけしか構文解析できない
func TestEM(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"*t -ext-*",
			"<p><em>t -ext-</em></p>\n",
		},
		{
			"**text**",
			"<p><strong>text</strong></p>\n",
		},
		{
			"***te xt***",
			"<p><strong><em>te xt</em></strong></p>\n",
		},
		{
			"text*text*",
			"<p>text<em>text</em></p>\n",
		},
		{
			"*text*text",
			"<p><em>text</em>text</p>\n",
		},
		{
			"text*text*text",
			"<p>text<em>text</em>text</p>\n",
		},
		{
			"*text* ***text***",
			"<p><em>text</em> <strong><em>text</em></strong></p>\n",
		},
		{
			"*text****text***",
			"<p><em>text</em><strong><em>text</em></strong></p>\n",
		},
		{
			"**text**`hoge`*huga*",
			"<p><strong>text</strong><code>hoge</code><em>huga</em></p>\n",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		document := p.ParseDocument()

		actual := document.String()
		if actual != tt.expected {
			t.Errorf("input=%q wong. expected=%q, got=%q", tt.input, tt.expected, actual)
		}
	}
}

// インラインコードのテスト
// ひとまず1行だけしか構文解析できない
func TestInlineCode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"`t ext`",
			"<p><code>t ext</code></p>\n",
		},
		{
			"text`text`text",
			"<p>text<code>text</code>text</p>\n",
		},
		{
			"`text`*hoge*",
			"<p><code>text</code><em>hoge</em></p>\n",
		},
		{
			"## h2`text`h2",
			"<h2>h2<code>text</code>h2</h2>\n",
		},
		{
			"- `a`",
			"<p>\n<ul>\n<li><code>a</code></li>\n</ul>\n</p>\n",
		},
		{
			"- ~~a~~",
			"<p>\n<ul>\n<li><s>a</s></li>\n</ul>\n</p>\n",
		},
		{
			"- a`b`c",
			"<p>\n<ul>\n<li>a<code>b</code>c</li>\n</ul>\n</p>\n",
		},
		{
			"- `a``b`c",
			"<p>\n<ul>\n<li><code>a</code><code>b</code>c</li>\n</ul>\n</p>\n",
		},
		{
			"- *a*`b`c",
			"<p>\n<ul>\n<li><em>a</em><code>b</code>c</li>\n</ul>\n</p>\n",
		},
		{
			"- a`b`*c*",
			"<p>\n<ul>\n<li>a<code>b</code><em>c</em></li>\n</ul>\n</p>\n",
		},
		{
			"- *a*`b`*c*",
			"<p>\n<ul>\n<li><em>a</em><code>b</code><em>c</em></li>\n</ul>\n</p>\n",
		},
		{
			"- `a`\n- b\n- `c`",
			"<p>\n<ul>\n<li><code>a</code></li>\n<li>b</li>\n<li><code>c</code></li>\n</ul>\n</p>\n",
		},
		{
			"`text``code`",
			"<p><code>text</code><code>code</code></p>\n",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		document := p.ParseDocument()

		actual := document.String()
		if actual != tt.expected {
			t.Errorf("input=%q wong. expected=%q, got=%q", tt.input, tt.expected, actual)
		}
	}
}

// Discリストのテスト
// まだ複数のリストはかけない
func TestDiscList(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"- 1text",
			`<p>
<ul>
<li>1text</li>
</ul>
</p>
`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		document := p.ParseDocument()

		actual := document.String()
		if actual != tt.expected {
			t.Errorf("input=%q wong. expected=%q, got=%q", tt.input, tt.expected, actual)
		}
	}
}

// 複数のDiscリストのテスト
func TestMultipleDiscList(t *testing.T) {
	input := `- a3
- b
- c4
`

	expected := `<p>
<ul>
<li>a3</li>
<li>b</li>
<li>c4</li>
</ul>
</p>
`

	l := lexer.New(input)
	p := New(l)
	document := p.ParseDocument()

	actual := document.String()
	if actual != expected {
		t.Errorf("input=%q wong. expected=%q, got=%q",
			input, expected, actual)
	}
}

// 複数のトークンの構文解析
func TestMultipeTokens(t *testing.T) {
	input := `# Heading

---

- a
- b
- c3

---

***te3 xt***
*1text*
**text2**



## Heading*2*`

	expected := `<h1>Heading</h1>
<hr>
<p>
<ul>
<li>a</li>
<li>b</li>
<li>c3</li>
</ul>
</p>
<hr>
<p><strong><em>te3 xt</em></strong><em>1text</em><strong>text2</strong></p>
<h2>Heading<em>2</em></h2>
`

	l := lexer.New(input)
	p := New(l)
	document := p.ParseDocument()

	actual := document.String()
	if actual != expected {
		t.Errorf("input=%q wong. expected=%q, got=%q",
			input, expected, actual)
	}
}

// 複数のトークンの構文解析
func TestMultipeTokens2(t *testing.T) {
	input := `
# *text*1MIDASHI1
## M2
*2text*
*text2*
~~strikethrough~~

## Heading*2* *text*
3text 999 hoge

---

### h3**!!!**
This is a text.
`

	expected := `<h1><em>text</em>1MIDASHI1</h1>
<h2>M2</h2>
<p><em>2text</em><em>text2</em><s>strikethrough</s></p>
<h2>Heading<em>2</em> <em>text</em></h2>
<p>3text 999 hoge</p>
<hr>
<h3>h3<strong>!!!</strong></h3>
<p>This is a text.</p>
`

	l := lexer.New(input)
	p := New(l)
	document := p.ParseDocument()

	actual := document.String()
	if actual != expected {
		t.Errorf("input=%q wong. expected=%q, got=%q",
			input, expected, actual)
	}
}

// コードブロックの構文解析
func TestCodeBlock(t *testing.T) {
	input := "```go\n"
	input += "func main() {\n"
	input += "    fmt.Printf(\"``Hello, world!``\")\n"
	input += "}\n"
	input += "```"

	expected := "<pre class=\"language-go\">\n"
	expected += "<code>\n"
	expected += "func main() {\n"
	expected += "    fmt.Printf(\"``Hello, world!``\")\n"
	expected += "}\n"
	expected += "</code>\n"
	expected += "</pre>\n"

	l := lexer.New(input)
	p := New(l)
	document := p.ParseDocument()

	actual := document.String()
	if actual != expected {
		t.Errorf("input=%q wong. expected=%q, got=%q",
			input, expected, actual)
	}
}

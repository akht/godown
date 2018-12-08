package ast

import (
	"bytes"
	"godown/token"
	"strconv"
)

type Node interface {
	TokenLiteral() string
	String() string
}

// ブロック要素
type Block interface {
	Node
	blockNode()
}

// インライン要素
type Inline interface {
	Node
	inlineNode()
}

// ASTのルートノード
type Document struct {
	Blocks []Block
}

func (d *Document) TokenLiteral() string {
	if len(d.Blocks) > 0 {
		return d.Blocks[0].TokenLiteral()
	} else {
		return ""
	}
}

func (d *Document) String() string {
	var out bytes.Buffer

	for _, b := range d.Blocks {
		out.WriteString(b.String())
	}

	return out.String()
}

// 見出し
type Heading struct {
	Token    token.Token
	Level    int
	Contents []Inline
}

func (h *Heading) blockNode()           {}
func (h *Heading) TokenLiteral() string { return h.Token.Literal }
func (h *Heading) String() string {
	var out bytes.Buffer

	htag := "h" + strconv.Itoa(h.Level)

	out.WriteString("<")
	out.WriteString(htag)
	out.WriteString(">")
	for _, l := range h.Contents {
		out.WriteString(l.String())
	}
	out.WriteString("</")
	out.WriteString(htag)
	out.WriteString(">")
	out.WriteString("\n")

	return out.String()
}

// Discリスト
type DiscList struct {
	Token token.Token
	Lists [][]Inline
}

func (d *DiscList) blockNode()           {}
func (d *DiscList) TokenLiteral() string { return d.Token.Literal }
func (d *DiscList) String() string {
	var out bytes.Buffer

	out.WriteString("<p>\n")
	out.WriteString("<ul>\n")

	for _, l := range d.Lists {
		out.WriteString("<li>")
		for _, l2 := range l {
			out.WriteString(l2.String())
		}
		out.WriteString("</li>\n")
	}

	out.WriteString("</ul>\n")
	out.WriteString("</p>")
	out.WriteString("\n")

	return out.String()
}

// パラグラフ
type Paragraph struct {
	Token    token.Token
	Contents []Inline
}

func (p *Paragraph) blockNode()           {}
func (p *Paragraph) TokenLiteral() string { return "" }
func (p *Paragraph) String() string {
	var out bytes.Buffer

	out.WriteString("<p>")

	for _, l := range p.Contents {
		out.WriteString(l.String())
	}

	out.WriteString("</p>")
	out.WriteString("\n")

	return out.String()
}

// コードブロック
type CodeBlock struct {
	Token    token.Token
	Lang     Inline
	Contents []Inline
}

func (c *CodeBlock) blockNode()           {}
func (c *CodeBlock) TokenLiteral() string { return "" }
func (c *CodeBlock) String() string {
	var out bytes.Buffer

	var lang string
	if c.Lang.String() == "" {
		lang = ""
	} else {
		lang = c.Lang.String()
	}

	out.WriteString("<pre class=\"language-" + lang + "\">")
	out.WriteString("\n")
	out.WriteString("<code>")
	out.WriteString("\n")

	for _, l := range c.Contents {
		out.WriteString(l.String())
	}

	out.WriteString("</code>")
	out.WriteString("\n")
	out.WriteString("</pre>")
	out.WriteString("\n")

	return out.String()
}

// 強調
type Emphasis struct {
	Token    token.Token
	Level    int
	Contents []Inline
}

func (e *Emphasis) inlineNode()          {}
func (e *Emphasis) TokenLiteral() string { return e.Token.Literal }
func (e *Emphasis) String() string {
	var out bytes.Buffer

	if e.Level == 1 {
		out.WriteString("<em>")
		for _, l := range e.Contents {
			out.WriteString(l.String())
		}
		out.WriteString("</em>")
	} else if e.Level == 2 {
		out.WriteString("<strong>")
		for _, l := range e.Contents {
			out.WriteString(l.String())
		}
		out.WriteString("</strong>")
	} else if e.Level == 3 {
		out.WriteString("<strong><em>")
		for _, l := range e.Contents {
			out.WriteString(l.String())
		}
		out.WriteString("</em></strong>")
	}

	return out.String()
}

// インラインコード
type InlineCode struct {
	Token    token.Token
	Contents []Inline
}

func (ic *InlineCode) inlineNode()          {}
func (ic *InlineCode) TokenLiteral() string { return ic.Token.Literal }
func (ic *InlineCode) String() string {
	var out bytes.Buffer

	out.WriteString("<code>")
	for _, l := range ic.Contents {
		out.WriteString(l.String())
	}
	out.WriteString("</code>")

	return out.String()
}

// 打ち消し
type Strikethrough struct {
	Token    token.Token
	Contents []Inline
}

func (s *Strikethrough) inlineNode()          {}
func (s *Strikethrough) TokenLiteral() string { return s.Token.Literal }
func (s *Strikethrough) String() string {
	var out bytes.Buffer

	out.WriteString("<s>")
	for _, l := range s.Contents {
		out.WriteString(l.String())
	}
	out.WriteString("</s>")

	return out.String()
}

// インラインテキスト
type Text struct {
	Token   token.Token
	Content string
}

func (t *Text) inlineNode()          {}
func (t *Text) TokenLiteral() string { return t.Token.Literal }
func (t *Text) String() string       { return t.Content }

// 水平線
type HorizontalRule struct{}

func (h *HorizontalRule) blockNode()           {}
func (h *HorizontalRule) TokenLiteral() string { return "" }
func (h *HorizontalRule) String() string       { return "<hr/>\n" }

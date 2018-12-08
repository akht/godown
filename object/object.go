package object

import (
	"bytes"
	"godown/decorator"
)

type ObjectType string

const (
	DOCUMENT_OBJ       = "DOCUMENT"
	HEADING_OBJ        = "HEADING"
	DISCLIST_OBJ       = "DISCLIST"
	CODEBLOCK_OBJ      = "CODEBLOCK"
	PARAGRAPH_OBJ      = "PARAGRAPH"
	HORIZONTALRULE_OBJ = "HORIZONTAL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// 文書全体のオブジェクト
type Document struct {
	Objects []Object
}

func (d *Document) Type() ObjectType { return DOCUMENT_OBJ }
func (d *Document) Inspect() string {
	var out bytes.Buffer

	for _, o := range d.Objects {
		out.WriteString(o.Inspect())
	}

	return out.String()
}
func (d *Document) Render() string {
	var out bytes.Buffer

	out.WriteString("<!DOCTYPE html>")
	out.WriteString("\n")
	out.WriteString("<html>")
	out.WriteString("\n")
	out.WriteString("<meta charset=\"utf-8\">")
	out.WriteString("\n")
	out.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")
	out.WriteString("\n")

	style(&out)
	out.WriteString("\n")

	out.WriteString("</head>")
	out.WriteString("\n")

	out.WriteString("<body for=\"html-export\" class=\"body\">")
	out.WriteString("\n")

	out.WriteString(d.Inspect())

	out.WriteString("</body>")
	out.WriteString("\n")
	out.WriteString("</html>")
	out.WriteString("\n")

	return out.String()
}

func style(out *bytes.Buffer) {
	out.WriteString("<style>\n")

	decorator.Deco(out)

	out.WriteString("</style>\n")
}

// 見出しを表現するオブジェクト
type Heading struct {
	Value string
}

func (h *Heading) Type() ObjectType { return HEADING_OBJ }
func (h *Heading) Inspect() string  { return h.Value }

// Discリストを表現するオブジェクト
type DiscList struct {
	Value string
}

func (dl *DiscList) Type() ObjectType { return DISCLIST_OBJ }
func (dl *DiscList) Inspect() string  { return dl.Value }

// コードブロックを表現するオブジェクト
type CodeBlock struct {
	Value string
}

func (c *CodeBlock) Type() ObjectType { return CODEBLOCK_OBJ }
func (c *CodeBlock) Inspect() string  { return c.Value }

// パラグラフを表現するオブジェクト
type Paragraph struct {
	Value string
}

func (p *Paragraph) Type() ObjectType { return PARAGRAPH_OBJ }
func (p *Paragraph) Inspect() string  { return p.Value }

// 水平線
type HorizontalRule struct{}

func (h *HorizontalRule) Type() ObjectType { return HORIZONTALRULE_OBJ }
func (h *HorizontalRule) Inspect() string  { return "<hr/>\n" }

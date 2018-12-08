package parser

import (
	"godown/ast"
	"godown/lexer"
	"godown/token"
)

// パーサ
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

// トークンを進める
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 現在のトークンが引数tと等しいか判定する
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次のトークンが引数tと等しいか判定する
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 現在のトークンがブロック要素かどうか
func (p *Parser) curTokenIsBlockNode() bool {
	return p.curTokenIs(token.IGETA) || p.curTokenIs(token.HYPHEN)
}

// 現在のトークンがインライン要素かどうか
func (p *Parser) curTokenIsInlineNode(context ast.Block) bool {

	isInlineNode := p.curTokenIs(token.ASTERISK) ||
		p.curTokenIs(token.BACKQUOTE) ||
		p.curTokenIs(token.TILDE) ||
		p.curTokenIs(token.TEXT) ||
		p.curTokenIs(token.SPACE)

	switch context.(type) {
	case *ast.Heading:
		// ヘッダーの文脈では、ハイフンもインラインテキストとして扱う
		return isInlineNode || p.curTokenIs(token.HYPHEN)
	}

	return isInlineNode
}

// peekTokenがMarkdownのキーワードかどうか
func (p *Parser) peekTokenIsKeywords() bool {
	return p.peekTokenIs(token.IGETA) ||
		p.peekTokenIs(token.HYPHEN) ||
		p.peekTokenIs(token.ASTERISK) ||
		p.peekTokenIs(token.BACKQUOTE) ||
		p.curTokenIs(token.TILDE) ||
		p.peekTokenIs(token.TEXT)
}

// アサーション関数
// peekTokenの型をチェックし、その型が期待された正しいものだった場合に限って
// nextToken()を読んでトークンを進める
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	return false
}

// Markdownドキュメントをパースする
func (p *Parser) ParseDocument() *ast.Document {
	document := &ast.Document{}
	document.Blocks = []ast.Block{}

	for !p.curTokenIs(token.EOF) {
		block := p.parseDocument()
		if block != nil {
			document.Blocks = append(document.Blocks, block)
		}

		if p.curTokenIs(token.SPACE) || p.curTokenIs(token.CR) {
			p.nextToken()
		}
	}

	return document
}

// トークンに応じてそれぞれの関数でパースする
func (p *Parser) parseDocument() ast.Block {
	switch p.curToken.Type {
	case token.IGETA:
		return p.parseHeading()
	case token.HYPHEN:
		return p.parseDiscList()
	case token.CR:
		return nil
	default:
		return p.parseParagraph()
	}
}

// 見出しの構文解析
func (p *Parser) parseHeading() *ast.Heading {
	block := &ast.Heading{Token: p.curToken}

	level := 1
	for p.expectPeek(token.IGETA) {
		level++
	}

	if !p.expectPeek(token.SPACE) {
		return nil
	}

	p.nextToken()

	block.Level = level
	inlineContent := p.parseInlineContent(block)
	block.Contents = inlineContent

	if p.curTokenIs(token.IGETA) || p.curTokenIs(token.HYPHEN) {
		return block
	}

	return block
}

// DISCリストの構文解析
func (p *Parser) parseDiscList() ast.Block {
	if p.curTokenIs(token.HYPHEN) && p.peekTokenIs(token.HYPHEN) {
		// --の場合は水平線としてパース
		return p.parseHorizontalRule()
	}

	DiscList := &ast.DiscList{}

	listtext := p.parseListItem(DiscList)
	DiscList.Lists = append(DiscList.Lists, listtext)

	for p.curTokenIs(token.HYPHEN) && !p.peekTokenIs(token.HYPHEN) {
		listtext := p.parseListItem(DiscList)
		DiscList.Lists = append(DiscList.Lists, listtext)
	}

	return DiscList
}

func (p *Parser) parseHorizontalRule() *ast.HorizontalRule {

	for p.curTokenIs(token.HYPHEN) {
		p.nextToken()
	}

	if p.curTokenIs(token.CR) {
		p.nextToken()
	}

	return &ast.HorizontalRule{}
}

// リストアイテムの構文解析
func (p *Parser) parseListItem(context ast.Block) []ast.Inline {

	if p.curTokenIs(token.HYPHEN) {
		p.nextToken()
	}
	if p.curTokenIs(token.SPACE) {
		p.nextToken()
	}

	listContent := p.parseInlineContent(context)

	if !p.curTokenIs(token.HYPHEN) {
		p.nextToken()
	}

	return listContent
}

// パラグラフのパース
func (p *Parser) parseParagraph() ast.Block {
	if p.curTokenIs(token.BACKQUOTE) && p.peekTokenIs(token.BACKQUOTE) {
		// ``の場合はコードブロックとしてパース
		return p.parseCodeBlock()
	}

	paragraph := &ast.Paragraph{}

	for !(p.curTokenIs(token.EOF) || p.curTokenIsBlockNode()) {
		inlineContent := p.parseInlineContent(paragraph)
		if inlineContent != nil {
			paragraph.Contents = append(paragraph.Contents, inlineContent...)
		}

		if p.curTokenIs(token.CR) || p.curTokenIs(token.TEXT) || p.curTokenIs(token.SPACE) {
			p.nextToken()
		}
	}

	return paragraph
}

// コードブロックのパース
func (p *Parser) parseCodeBlock() ast.Block {
	codeBlock := &ast.CodeBlock{}

	for p.curTokenIs(token.BACKQUOTE) {
		p.nextToken()
	}

	if !p.curTokenIs(token.CR) {
		codeBlock.Lang = p.parseInlineText()
		p.nextToken()
		p.nextToken()
	}

	for !p.curTokenIs(token.BACKQUOTE) {
		codeBlock.Contents = p.parseCodeBlockContent()
	}

	for {
		var tmp []ast.Inline
		count := 0
		for p.curTokenIs(token.BACKQUOTE) {
			text := &ast.Text{Token: p.curToken, Content: p.curToken.Literal}
			tmp = append(tmp, text)
			p.nextToken()
			count++
		}

		if count == 3 {
			break
		} else {
			codeBlock.Contents = append(codeBlock.Contents, tmp...)
			for !p.curTokenIs(token.BACKQUOTE) {
				codeBlock.Contents = append(codeBlock.Contents, p.parseCodeBlockContent()...)
			}
		}
	}

	return codeBlock
}

func (p *Parser) parseCodeBlockContent() []ast.Inline {
	var contents []ast.Inline

	for !p.curTokenIs(token.BACKQUOTE) {
		contents = append(contents, p.parseInlineText())
		p.nextToken()
	}

	return contents
}

// インライン要素の構文解析
func (p *Parser) parseInlineContent(context ast.Block) []ast.Inline {
	var inlineContents []ast.Inline

	for p.curTokenIsInlineNode(context) {
		var inlineContent ast.Inline

		switch p.curToken.Type {
		case token.ASTERISK:
			inlineContent = p.parseInlineEmphasis()
		case token.BACKQUOTE:
			inlineContent = p.parseInlineCode()
		case token.TILDE:
			inlineContent = p.parseInlineStrikethrough()
		default:
			inlineContent = p.parseInlineText()
			p.nextToken()
		}

		if inlineContent != nil {
			inlineContents = append(inlineContents, inlineContent)
		}

		switch context.(type) {
		case *ast.Heading:
			// ヘッダーの文脈では、改行が来たらブロック終了
			if p.curTokenIs(token.CR) {
				p.nextToken()
				return inlineContents
			}
		}

		if !p.curTokenIsInlineNode(context) {
			p.nextToken()
		}

		if p.curTokenIs(token.CR) && p.peekTokenIsKeywords() {
			return inlineContents
		}
	}

	return inlineContents
}

// 強調の構文解析
func (p *Parser) parseInlineEmphasis() *ast.Emphasis {
	emphasis := &ast.Emphasis{Token: p.curToken}

	level := 1
	for p.expectPeek(token.ASTERISK) {
		level++
	}

	p.nextToken()

	emphasis.Level = level

	for !p.curTokenIs(token.ASTERISK) {
		emphasis.Contents = append(emphasis.Contents, p.parseInlineText())
		p.nextToken()
	}

	for i := 0; i < level; i++ {
		p.nextToken()
	}

	return emphasis
}

// 打ち消しの構文解析
func (p *Parser) parseInlineStrikethrough() *ast.Strikethrough {
	strikethrough := &ast.Strikethrough{Token: p.curToken}

	count := 1
	for p.expectPeek(token.TILDE) {
		count++
	}

	p.nextToken()

	for !p.curTokenIs(token.TILDE) {
		strikethrough.Contents = append(strikethrough.Contents, p.parseInlineText())
		p.nextToken()
	}

	for i := 0; i < count; i++ {
		p.nextToken()
	}

	return strikethrough
}

// インラインコードの構文解析
func (p *Parser) parseInlineCode() *ast.InlineCode {
	inlineCode := &ast.InlineCode{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.BACKQUOTE) {
		inlineCode.Contents = append(inlineCode.Contents, p.parseInlineText())
		p.nextToken()
	}

	p.nextToken()

	return inlineCode
}

// インラインテキストの構文解析
func (p *Parser) parseInlineText() ast.Inline {
	return &ast.Text{Token: p.curToken, Content: p.curToken.Literal}
}

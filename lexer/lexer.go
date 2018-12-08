package lexer

import "godown/token"

// レキサー
type Lexer struct {
	input        string // 入力されたMarkdown文書
	position     int    // 入力における現在の位置
	readPosition int    // 現在の文字の次の文字
	ch           byte
	ch_debug     string // デバッグ用のch
}

// Markdown文書からレキサーを生成
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipCarriageReturn()

	switch l.ch {
	case '#':
		tok = newToken(token.IGETA, l.ch)
	case ' ':
		tok = newToken(token.SPACE, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '-':
		tok = newToken(token.HYPHEN, l.ch)
	case '`':
		tok = newToken(token.BACKQUOTE, l.ch)
	case '~':
		tok = newToken(token.TILDE, l.ch)
	case '\n':
		tok = newToken(token.CR, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		ch := l.ch
		literal := string(ch) + l.readText()
		tok.Type = token.TEXT
		tok.Literal = literal
		return tok
	}

	l.readChar()
	return tok
}

// 次の１文字を読み込んで、inputの現在位置を進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.ch_debug = string(l.ch)
	} else {
		l.ch = l.input[l.readPosition]
		l.ch_debug = string(l.ch)
	}
	l.position = l.readPosition
	l.readPosition++
}

// 次の１文字をpeekする
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 改行を読み飛ばす
func (l *Lexer) skipCarriageReturn() {
	// for l.ch == '\n' || l.ch == '\r' {
	for l.ch == '\r' {
		l.readChar()
	}
}

// 文字列を読み取る(改行文字か入力の最後に至るまでreadChar()を呼ぶ)
func (l *Lexer) readText() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == '\r' || l.ch == 0 || l.ch == ' ' || l.ch == '#' || l.ch == '*' || l.ch == '-' || l.ch == '`' || l.ch == '~' {
			break
		}
	}

	return l.input[position:l.position]
}

// 文字が数字かどうか判定する
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 数字を読み込んで、非数字に達するまで字句解析器の位置を進める
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

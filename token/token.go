package token

// トークンの種類を区別するための型
type TokenType string

// トークン
type Token struct {
	Type    TokenType // トークンの種類を区別する
	Literal string    // トークンのリテラル表現
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	CR      = "CR"

	IGETA  = "#"
	HYPHEN = "-"

	ASTERISK  = "*"
	BACKQUOTE = "`"
	TILDE     = "~"
	TEXT      = "TEXT" // 文字列
	SPACE     = " "
	// INT       = "INT"  // 数字
	// DOT       = "."
)

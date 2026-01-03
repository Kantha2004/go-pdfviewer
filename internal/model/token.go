package model

// TokenType represents the category of a lexical token.
type TokenType int

// TokenType values.
const (
	TokEOF TokenType = iota
	TokNumber
	TokName
	TokString
	TokHexaString
	TokArrayStart
	TokArrayEnd
	TokDictStart
	TokDictEnd
	TokKeyword
)

// Delimiter and punctuation bytes used by the lexer.
const (
	OpenParen     byte = '('
	CloseParen    byte = ')'
	LessThan      byte = '<'
	GreaterThan   byte = '>'
	OpenLBracket  byte = '['
	CloseLBracket byte = ']'
	OpenBrace     byte = '{'
	CloseBrace    byte = '}'
	Solidus       byte = '/'
)

const (
	Percent byte = '%'
)

const (
	Minus   byte = '-'
	Plus    byte = '+'
	Decimal byte = '.'
)

// Token represents a single lexical token produced by the lexer.
type Token struct {
	Type  TokenType
	Value string
}

func (t TokenType) String() string {
	switch t {
	case TokEOF:
		return "EOF"
	case TokNumber:
		return "Number"
	case TokName:
		return "Name"
	case TokString:
		return "String"
	case TokHexaString:
		return "HexaString"
	case TokArrayStart:
		return "ArrayStart"
	case TokArrayEnd:
		return "ArrayEnd"
	case TokDictStart:
		return "DictStart"
	case TokDictEnd:
		return "DictEnd"
	case TokKeyword:
		return "Keyword"
	default:
		return "UnknownTokenType"
	}
}

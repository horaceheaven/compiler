package lexer

import (
	"github.com/horaceheaven/compiler/token"
)

type Lexer struct {
	position int

	readPosition int

	ch rune

	characters []rune

	prevToken token.Token
}

func New(input string) *Lexer {
	l := &Lexer{characters: []rune(input)}
	l.readChar()
	return l
}

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

// Package parser is used to parse input-programs written in monkey
// and convert them to an abstract-syntax tree.
package parser

import (
	"github.com/horaceheaven/compiler/ast"
	"github.com/horaceheaven/compiler/lexer"
	"github.com/horaceheaven/compiler/token"
)

// prefix Parse function
// infix parse function
// postfix parse function
type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func() ast.Expression
)

// precedence order
const (
	_ int = iota
	LOWEST
	COND         // OR or AND
	ASSIGN       // =
	TERNARY      // ? :
	EQUALS       // == or !=
	REGEXP_MATCH // !~ ~=
	LESSGREATER  // > or <
	SUM          // + or -
	PRODUCT      // * or /
	POWER        // **
	MOD          // %
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	DOTDOT       // ..
	INDEX        // array[index], map[key]
	HIGHEST
)

// each token precedence
var precedences = map[token.Type]int{
	token.QUESTION:     TERNARY,
	token.ASSIGN:       ASSIGN,
	token.DOTDOT:       DOTDOT,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.LT_EQUALS:    LESSGREATER,
	token.GT:           LESSGREATER,
	token.GT_EQUALS:    LESSGREATER,
	token.CONTAINS:     REGEXP_MATCH,
	token.NOT_CONTAINS: REGEXP_MATCH,

	token.PLUS:            SUM,
	token.PLUS_EQUALS:     SUM,
	token.MINUS:           SUM,
	token.MINUS_EQUALS:    SUM,
	token.SLASH:           PRODUCT,
	token.SLASH_EQUALS:    PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.ASTERISK_EQUALS: PRODUCT,
	token.POW:             POWER,
	token.MOD:             MOD,
	token.AND:             COND,
	token.OR:              COND,
	token.LPAREN:          CALL,
	token.PERIOD:          CALL,
	token.LBRACKET:        INDEX,
}

// Parser object
type Parser struct {
	// l is our lexer
	l *lexer.Lexer

	// prevToken holds the previous token from our lexer.
	// (used for "++" + "--")
	prevToken token.Token

	// curToken holds the current token from our lexer.
	curToken token.Token

	// peekToken holds the next token which will come from the lexer.
	peekToken token.Token

	// errors holds parsing-errors.
	errors []string

	// prefixParseFns holds a map of parsing methods for
	// prefix-based syntax.
	prefixParseFns map[token.Type]prefixParseFn

	// infixParseFns holds a map of parsing methods for
	// infix-based syntax.
	infixParseFns map[token.Type]infixParseFn

	// postfixParseFns holds a map of parsing methods for
	// postfix-based syntax.
	postfixParseFns map[token.Type]postfixParseFn

	// are we inside a ternary expression?
	//
	// Nested ternary expressions are illegal :)
	tern bool
}

// New returns our new parser-object.
func New(l *lexer.Lexer) *Parser {

	// Create the parser, and prime the pump
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	// Register prefix-functions
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.BACKTICK, p.parseBacktickLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.DEFINE_FUNCTION, p.parseFunctionDefinition)
	p.registerPrefix(token.EOF, p.parsingBroken)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.FOR, p.parseForLoopExpression)
	p.registerPrefix(token.FOREACH, p.parseForEach)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.ILLEGAL, p.parsingBroken)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NULL, p.parseNull)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.SWITCH, p.parseSwitchStatement)
	p.registerPrefix(token.TRUE, p.parseBoolean)

	// Register infix functions
	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.CONTAINS, p.parseInfixExpression)
	p.registerInfix(token.DOTDOT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.NOT_CONTAINS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.PERIOD, p.parseMethodCallExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.PLUS_EQUALS, p.parseAssignExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.SLASH_EQUALS, p.parseAssignExpression)

	// Register postfix functions.
	p.postfixParseFns = make(map[token.Type]postfixParseFn)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)

	// All done
	return p
}

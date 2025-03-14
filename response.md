I'll help you create a lexer and parser for CCL in Go. Let's break this down into steps.

First, let's define the token types for the lexer:

```go
// token/token.go
package token

type TokenType string

const (
    ILLEGAL   TokenType = "ILLEGAL"
    EOF       TokenType = "EOF"
    
    // Identifiers + literals
    IDENT     TokenType = "IDENT"
    INT       TokenType = "INT"
    STRING    TokenType = "STRING"
    
    // Operators and delimiters
    SEMICOLON TokenType = ";"
    COLON     TokenType = ";"
    LBRACE    TokenType = "{"
    RBRACE    TokenType = "}"
    LBRACKET  TokenType = "["
    RBRACKET  TokenType = "]"
    LPAREN    TokenType = "("
    RPAREN    TokenType = ")"
    
    // Keywords
    MODEL     TokenType = "MODEL"
    
    // Attributes
    ATTRIBUTE TokenType = "#"
)

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

Now, let's create a basic lexer:

```go
// lexer/lexer.go
package lexer

import "your/project/token"

type Lexer struct {
    input        string
    position     int  // current position in input (points to current char)
    readPosition int  // current reading position in input (after current char)
    ch           byte // current char under examination
    line         int
    column       int
}

func New(input string) *Lexer {
    l := &Lexer{
        input:    input,
        line:     1,
        column:   1,
    }
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
    l.column += 1
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    l.skipWhitespace()

    switch l.ch {
    case '{':
        tok = newToken(token.LBRACE, l.ch, l.line, l.column)
    case '}':
        tok = newToken(token.RBRACE, l.ch, l.line, l.column)
    case '[':
        tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
    case ']':
        tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
    case '(':
        tok = newToken(token.LPAREN, l.ch, l.line, l.column)
    case ')':
        tok = newToken(token.RPAREN, l.ch, l.line, l.column)
    case ':':
        tok = newToken(token.COLON, l.ch, l.line, l.column)
    case ';':
        tok = newToken(token.SEMICOLON, l.ch, l.line, l.column)
    case '#':
        tok = newToken(token.ATTRIBUTE, l.ch, l.line, l.column)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            tok.Line = l.line
            tok.Column = l.column
            return tok
        } else if isDigit(l.ch) {
            tok.Literal = l.readNumber()
            tok.Type = token.INT
            tok.Line = l.line
            tok.Column = l.column
            return tok
        } else {
            tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        if l.ch == '\n' {
            l.line += 1
            l.column = 0
        }
        l.readChar()
    }
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line, column int) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}
```

This is a basic start. For the parser, we'll need to define AST nodes:

```go
// ast/ast.go
package ast

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

type ModelStatement struct {
    Token      token.Token
    Name       *Identifier
    Attributes []*Attribute
    Fields     []*Field
}

type Attribute struct {
    Token token.Token
    Name  string
    Value string
}

type Field struct {
    Token token.Token
    Name  *Identifier
    Type  *Identifier
    IsArray bool
}

type Identifier struct {
    Token token.Token
    Value string
}
```

Would you like me to continue with the parser implementation? The parser will be more complex and will need to handle:
1. Attribute parsing
2. Model declarations
3. Field declarations
4. Type declarations
5. Error handling

Let me know if you want me to proceed with any specific part!
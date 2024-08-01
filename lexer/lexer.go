package lexer

import (
	// "strconv"
	"errors"
	"strconv"
	"strings"

	"github.com/jf550-kent/jsgo/token"
)

// Lexer tokenization the source text for the language
type Lexer struct {
	src          []byte // source text for tokenization
	position     int    // current position at [Lexer.src]
	nextPosition int    // next position to be lex at [Lexer.src]

	line int // the current line at the source text
	col  int // the current column at the source text

	ch byte // the current byte at [Lexer.src]
}

// New return a *Lexer
func New(byt []byte) *Lexer {
	l := &Lexer{src: byt, line: 1}
	l.next()
	return l
}

// Lex return the next token in the [*Lexer]
func (l *Lexer) Lex() (token.Token, error) {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '+':
		pos := l.currentPos()
		tok = newToken(token.ADD, "+", pos, pos)
	case '-':
		pos := l.currentPos()
		tok = newToken(token.MINUS, "-", pos, pos)
	case '*':
		pos := l.currentPos()
		tok = newToken(token.MUL, "*", pos, pos)
	case '/':
		pos := l.currentPos()
		tok = newToken(token.DIVIDE, "/", pos, pos)
	case ',':
		pos := l.currentPos()
		tok = newToken(token.COMMA, ",", pos, pos)
	case '.':
		pos := l.currentPos()
		tok = newToken(token.DOT, ".", pos, pos)
	case ':':
		pos := l.currentPos()
		tok = newToken(token.COLON, ":", pos, pos)
	case ';':
		pos := l.currentPos()
		tok = newToken(token.SEMICOLON, ";", pos, pos)
	case '(':
		pos := l.currentPos()
		tok = newToken(token.LPAREN, "(", pos, pos)
	case ')':
		pos := l.currentPos()
		tok = newToken(token.RPAREN, ")", pos, pos)
	case '{':
		pos := l.currentPos()
		tok = newToken(token.LBRACE, "{", pos, pos)
	case '}':
		pos := l.currentPos()
		tok = newToken(token.RBRACE, "}", pos, pos)
	case '>':
		pos := l.currentPos()
		tok = newToken(token.GTR, ">", pos, pos)
	case '<':
		pos := l.currentPos()
		tok = newToken(token.LSS, "<", pos, pos)
	case '"':
		tok = l.readString()
	case '!':
		start := l.currentPos()
		if l.peekByte() == '=' {
			l.next()
			end := l.currentPos()
			tok = newToken(token.NOT_EQUAL, "!=", start, end)
			break
		}
		tok = newToken(token.BANG, "!", start, start)
	case '=':
		start := l.currentPos()
		if l.peekByte() == '=' {
			l.next()
			end := l.currentPos()
			tok = newToken(token.EQUAL, "==", start, end)
			break
		}
		tok = newToken(token.ASSIGN, "=", start, start)
	case 0:
		return newToken(token.EOF, "EOF", l.currentPos(), l.currentPos()), nil
	default:
		if l.isLetter() {
			start := l.currentPos()
			s, end := l.getLetter()
			if ty, ok := token.Keyword(s); ok {
				return newToken(ty, s, start, end), nil
			}
			return newToken(token.IDENT, s, start, end), nil
		}
		if l.isDigit() {
			return l.getDigitToken()
		}
		return newToken(token.ILLEGAL, "ILLEGAL", l.currentPos(), l.currentPos()), errors.New("ILLEGAL token")
	}
	l.next()
	return tok, nil
}

// getDigitToken returns either [token.Token.NUMBER] or [token.Token.FLOAT]
// with its corresponding literal
func (l *Lexer) getDigitToken() (token.Token, error) {
	var digit strings.Builder
	start := l.currentPos()
	var end token.Pos
	var err error
	hasDot := false
	for {
		if l.isDigit() {
			digit.WriteByte(l.ch)
			end = l.currentPos()
			l.next()
			continue
		}
		if l.ch == '.' {
			if hasDot {
				err = errors.New("digit formatted incorrect at " + strconv.Itoa(l.position))
			}
			hasDot = true
			digit.WriteByte(l.ch)
			end = l.currentPos()
			l.next()
			continue
		}
		break
	}
	if hasDot {
		return newToken(token.FLOAT, digit.String(), start, end), err
	}
	return newToken(token.NUMBER, digit.String(), start, end), err
}

// getLetter return the whole letter with the position
func (l *Lexer) getLetter() (string, token.Pos) {
	var letter strings.Builder
	var end token.Pos
	for l.isLetter() {
		letter.WriteByte(l.ch)
		end = l.currentPos()
		l.next()
	}
	return letter.String(), end
}

// next moves the current position of the char in [Lexer.data] to the next one
// it will the [Lexer.ch] to 0 when [Lexer.position] is at the last byte of the [Lexer.src]
func (l *Lexer) next() {
	if l.nextPosition == len(l.src) {
		l.ch = 0
		l.position = l.nextPosition
		return
	}
	l.position = l.nextPosition
	l.ch = l.src[l.position]
	switch l.ch {
	case '\n':
		l.col = 0
		l.line++
	default:
		l.col++
	}
	l.nextPosition++
}

// peekByte returns the next byte in [Lexer.src] without
// moviing the [Lexer.position] and [Lexer.NextPostion]
// use this function to look ahead of the [Lexer.src]
func (l *Lexer) peekByte() byte {
	if l.nextPosition == len(l.src) {
		return 0 // 0 represent End of file
	}
	return l.src[l.nextPosition]
}

func (l *Lexer) readString() token.Token {
	start := l.position
	strTok := token.Token{TokenType: token.STRING, Start: l.currentPos(), }
	l.next()
	for {
		l.next()

		if l.ch == '"' || l.ch == 0 {
			prvPos := l.position -1 // this logic needs to be improve, we cannot just minus one to get the last position
			if l.src[prvPos] != '\\' {
				break
			}
		}
	}
	strTok.End = l.currentPos()
	strTok.Literal = string(l.src[start:l.nextPosition])
	return strTok
}

// skipWhitespace skips all the current whitespace in [Lexer.src]
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.next()
	}
}

func (l *Lexer) currentPos() token.Pos {
	return token.Pos{Line: l.line, Col: l.col}
}

func newToken(typ token.TokenType, literal string, start, end token.Pos) token.Token {
	return token.Token{TokenType: typ, Literal: literal, Start: start, End: end}
}

func (l *Lexer) isLetter() bool {
	return 'a' <= l.ch && l.ch <= 'z' || l.ch == '_' || 'A' <= l.ch && l.ch <= 'Z'
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9'
}

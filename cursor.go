package goindex

import (
	"go/scanner"
	"go/token"
)

// NewCursor returns a cursor for the given scanner.
// The scanner should not be used outside of the cursor.
func NewCursor(s *scanner.Scanner) *Cursor {
	return &Cursor{s: s}
}

type Cursor struct {
	s   *scanner.Scanner
	tok token.Token
	pos token.Pos
	lit string

	paren int
	brace int
}

// Next returns true  until token.EOF is found
func (c *Cursor) Next() bool {
	pos, tok, lit := c.s.Scan()
	c.tok = tok
	c.pos = pos
	c.lit = lit
	c.feed(tok)

	return tok != token.EOF
}

func (c *Cursor) At(file *token.File) int {
	end := c.Pos()
	if end == 0 {
		return file.Size() - 1
	}
	return file.Offset(c.Pos())
}

func (c *Cursor) Pos() token.Pos     { return c.pos }
func (c *Cursor) Token() token.Token { return c.tok }
func (c *Cursor) Lit() string        { return c.lit }

func (c *Cursor) InsideParen() bool { return c.paren > 0 }
func (c *Cursor) InsideBrace() bool { return c.brace > 0 }

func (c *Cursor) scanSignature() token.Pos {
	for c.Next() {
		if !c.InsideParen() {
			break
		}
	}
	return c.Pos()
}

func (c *Cursor) scanParenBlock() token.Pos {
	for c.Next() {
		if c.Token() == token.SEMICOLON && !c.InsideParen() {
			break
		}
	}
	return c.Pos()
}

func (c *Cursor) scanBlockStart() token.Pos {
	for c.Next() {
		if c.Token() == token.LBRACE && !c.InsideParen() {
			break
		}
	}
	return c.Pos()
}

func (c *Cursor) scanBlockEnd() token.Pos {
	for c.Next() {
		if !c.InsideParen() && !c.InsideBrace() {
			break
		}
	}
	return c.Pos()
}

func (c *Cursor) feed(tok token.Token) {
	switch tok {
	case token.LPAREN:
		c.paren++
	case token.RPAREN:
		c.paren--
	case token.LBRACE:
		c.brace++
	case token.RBRACE:
		c.brace--
	}
}

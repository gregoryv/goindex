package goindex

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func Index(src []byte) []Section {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	c := NewCursor(&s)

	sections := make([]Section, 0)

	var from int
	for c.Next() {
		pos := c.Pos()
		if c.Token() == token.COMMENT {
			from = file.Offset(pos) // and position to include in func blocks
			l := len(c.Lit())
			prefix := string(src[from+l+1 : from+l+5]) // if related it's either func or type
			//fmt.Printf("l=%v %q\n", l, prefix)
			if prefix != "func" {
				from = -1
			}
		}
		switch c.Token() {
		case token.IMPORT:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}
			c.scanSignature()
			to := c.At(file) + 1
			sections = append(sections, &importSect{
				span: span{from: from, to: to},
			})

		case token.TYPE:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}
			c.Next()
			name := c.Lit()
			c.Next()
			variant := c.Lit()
			c.scanBlockStart()
			c.scanBlockEnd()
			to := c.At(file) + 1
			sections = append(sections, &typeSect{
				span: span{
					from: from,
					to:   to,
				},
				name:    name,
				variant: variant,
			})
		case token.FUNC:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}
			// todo find name of func
			c.Next()
			var name string
			if c.Token() == token.LPAREN {
				// skip the receiver
				for c.Next() {
					if c.Token() == token.RPAREN {
						break
					}
				}
				c.Next()
				if c.Token() == token.IDENT {
					name = c.Lit()
				}

			} else {
				if c.Token() == token.IDENT {
					name = c.Lit()
				}
			}
			c.scanSignature()
			c.scanBlockStart()
			c.scanBlockEnd()
			to := c.At(file) + 1
			sections = append(sections, &funcSect{
				span: span{
					from: from,
					to:   to,
				},
				name: name,
			})
		}
		if c.Token() != token.COMMENT {
			from = -1
		}
	}
	// insert missing sections
	res := make([]Section, 0)
	var to int
	for _, s := range sections {
		if to < s.From() {
			res = append(res, newOtherSect(to, s.From()))
			to = s.To()
		}
		res = append(res, s)
	}
	last := res[len(res)-1]
	if last.To() != len(src) {
		res = append(res, newOtherSect(last.To(), len(src)))
	}

	return res
}

// ----------------------------------------

type Section interface {
	From() int
	To() int
	String() string
}

type importSect struct {
	span
}

func (me *importSect) String() string { return "import" }

type typeSect struct {
	span
	name    string
	variant string // struct or interface
}

func (me *typeSect) String() string {
	return fmt.Sprintf("type %s %s", me.name, me.variant)
}

// ----------------------------------------

type funcSect struct {
	span

	name string
}

func (me *funcSect) String() string { return "func " + me.name }

// ----------------------------------------

func newOtherSect(from, to int) *otherSect {
	return &otherSect{span: span{from: from, to: to}}
}

type otherSect struct {
	span
}

func (me *otherSect) String() string { return "?" }

// ----------------------------------------

type span struct {
	from, to int
}

func (me *span) From() int { return me.from }
func (me *span) To() int   { return me.to }

// ----------------------------------------

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
func (me *Cursor) Next() bool {
	pos, tok, lit := me.s.Scan()
	me.tok = tok
	me.pos = pos
	me.lit = lit
	me.feed(tok)

	return tok != token.EOF
}

func (me *Cursor) At(file *token.File) int {
	end := me.Pos()
	if end == 0 {
		return file.Size() - 1
	}
	return file.Offset(me.Pos())
}

func (me *Cursor) Pos() token.Pos     { return me.pos }
func (me *Cursor) Token() token.Token { return me.tok }
func (me *Cursor) Lit() string        { return me.lit }

func (c *Cursor) scanSignature() token.Pos {
	for c.Next() {
		if !c.InsideParen() {
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

func (me *Cursor) feed(tok token.Token) {
	switch tok {
	case token.LPAREN:
		me.paren++
	case token.RPAREN:
		me.paren--
	case token.LBRACE:
		me.brace++
	case token.RBRACE:
		me.brace--
	}
}

func (me *Cursor) InsideParen() bool { return me.paren > 0 }
func (me *Cursor) InsideBrace() bool { return me.brace > 0 }

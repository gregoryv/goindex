package gosort

import (
	"bytes"
	"go/scanner"
	"go/token"
	"os"
	"strings"
	"testing"
)

func TestIndex_usingScanner(t *testing.T) {
	src, err := os.ReadFile("testdata/complex.go")
	if err != nil {
		t.Fatal(err)
	}

	sections := Index(src)

	for _, s := range sections {
		var buf bytes.Buffer
		buf.Write(src[s.From():s.To()])
		got := buf.String()
		if strings.Contains(got, "Decoupled comment") {
			t.Log(got)
			t.Error("contains unrelated comment")
		}
	}

}

func Index(src []byte) []Section {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	c := NewCursor(&s)

	sections := make([]Section, 0)
	// todo fix the comment relation
	var comment string
	var from int
	for c.Next() {
		pos := c.Pos()
		if c.Token() == token.COMMENT {
			comment = c.Lit()       // save comment
			from = file.Offset(pos) // and position to include in func blocks
		}
		if c.Token() == token.FUNC {
			if comment == "" { // comment close by
				from = file.Offset(pos)
			}
			c.scanSignature()
			c.scanBlockStart()
			c.scanBlockEnd()
			end := c.Pos()
			if end == 0 {
				panic("missing block end")
			}
			to := file.Offset(end) + 1
			sections = append(sections, &funcSect{
				span: span{
					from: from,
					to:   to,
				},
			})
		}
		if c.Token() != token.COMMENT {
			comment = ""
		}
	}
	return sections
}

// ----------------------------------------

type Section interface {
	From() int
	To() int
}

type funcSect struct {
	span
}

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

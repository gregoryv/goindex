package goindex

import (
	"bytes"
	"go/scanner"
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func TestIndex(t *testing.T) {
	src, err := os.ReadFile("testdata/complex.go")
	if err != nil {
		t.Fatal(err)
	}
	sections := Index(src)

	t.Run("decoupled comment", func(t *testing.T) {
		for _, s := range sections {
			var buf bytes.Buffer
			buf.Write(src[s.From():s.To()])
			got := buf.String()
			if strings.Contains(got, "Decoupled comment") {
				if strings.Contains(got, "func") || strings.Contains(got, "type") {
					t.Log(got)
					t.Error("contains unrelated comment")
				}
			}
		}
	})

	t.Run("related comment", func(t *testing.T) {
		var buf bytes.Buffer
		for _, s := range sections {
			buf.Write(src[s.From():s.To()])
		}
		got := buf.String()
		if !strings.Contains(got, "Func comment") {
			t.Log(got)
			t.Error("missing related comment")
		}
	})

	t.Run("starts from 0", func(t *testing.T) {
		if got := sections[0].From(); got != 0 {
			t.Error(got)
		}
	})

	t.Run("ends with src len", func(t *testing.T) {
		exp := len(src)
		if got := sections[len(sections)-1].To(); got != exp {
			t.Error(got, "expected", exp)
		}
	})

	t.Run("sections are complete", func(t *testing.T) {
		for i, s := range sections[:len(sections)-1] {
			next := sections[i+1]
			t.Log(s.To(), next.From())
			if s.To() != next.From() {
				t.Errorf("missing section between %v and %v", i, i+1)
			}
		}
	})

	t.Run("equals", func(t *testing.T) {
		var buf bytes.Buffer
		for _, s := range sections {
			buf.Write(src[s.From():s.To()])
		}
		got, exp := buf.String(), string(src)
		golden.AssertEquals(t, got, exp)
	})
}

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
		if c.Token() == token.FUNC {
			if from == -1 { // no related comment
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
			from = -1
		}
	}
	// insert the first section if it's unspecified
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
	Decl() string
}

type funcSect struct {
	span
}

func (me *funcSect) Decl() string { return "func" }

type span struct {
	from, to int
}

func (me *span) From() int { return me.from }
func (me *span) To() int   { return me.to }

func newOtherSect(from, to int) *otherSect {
	return &otherSect{span: span{from: from, to: to}}
}

type otherSect struct {
	span
}

func (me *otherSect) Decl() string { return "other" }

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

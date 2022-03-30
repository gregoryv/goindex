package gosort

import (
	"go/scanner"
	"go/token"
	"os"
	"testing"
)

func TestIndex_usingScanner(t *testing.T) {
	src, err := os.ReadFile("testdata/complex.go")
	if err != nil {
		t.Fatal(err)
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.FUNC {
			from := file.Offset(pos)
			_ = scanSignature(&s)
			_ = scanBlockStart(&s)
			end := scanBlockEnd(&s)
			if end == 0 {
				t.Fatal("missing block end")
			}
			to := file.Offset(end) + 1
			t.Error(string(src[from:to]))
		}
	}
}

func scanSignature(s *scanner.Scanner) token.Pos {
	var count int
	var position token.Pos
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.LPAREN {
			count++
			continue
		}
		if tok == token.RPAREN {
			count--
			if count == 0 {
				position = pos
				break
			}
		}
	}
	return position
}

func scanBlockStart(s *scanner.Scanner) token.Pos {
	var inside bool
	var position token.Pos
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.LPAREN {
			inside = true
			continue
		}
		if tok == token.RPAREN {
			inside = false
			continue
		}

		if tok == token.LBRACE && !inside {
			position = pos
			break
		}
	}
	return position
}

func scanBlockEnd(s *scanner.Scanner) token.Pos {
	var c Cursor
	var position token.Pos

	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			break
		}
		c.Feed(tok)
		if !c.InsideParen() && !c.InsideBrace() {
			position = pos
			break
		}
	}
	return position
}

type Cursor struct {
	paren int
	brace int
}

func (me *Cursor) Feed(tok token.Token) {
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

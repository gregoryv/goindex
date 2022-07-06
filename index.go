package goindex

import (
	"go/scanner"
	"go/token"
)

// Index parses the given go source into sections. A section can be
// - comment
// - import
// - type
// - func
// If a comment is coupled to e.g. a func it's included in that section.
func Index(src []byte) []Section {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	c := NewCursor(&s)

	sections := make([]Section, 0)

	var from int
	var lastTok token.Token
	for c.Next() {
		pos := c.Pos()
		tok := c.Token()
		switch tok {
		case token.COMMENT:
			if lastTok == tok {
				continue // multiline comment
			}
			from = file.Offset(pos) // and position to include in func blocks
			l := len(c.Lit())
			// check if this comment is related to a func or type
			end := from + l + 5 // newline + 4 bytes = 5
			if end >= len(src) {
				end = len(src)
			}
			v := string(src[from+l+1 : end]) // if related it's either func or type
			switch v {
			case "func", "type":
			default:
				from = -1
			}
			//fmt.Printf("l=%v %q\n", l, v)

		case token.IMPORT:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}
			c.scanParenBlock()
			to := c.At(file) + 1
			sections = append(sections, newImport(from, to))

		case token.TYPE:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}
			c.scanBlockStart()
			label := string(src[file.Offset(pos):file.Offset(c.Pos())])
			c.scanBlockEnd()
			to := c.At(file) + 1
			sections = append(sections, newSection(from, to, label))

		case token.FUNC:
			if from == -1 { // no related comment
				from = file.Offset(pos)
			}

			// Fixme: func extra(), ie. no body
			// See https://go.dev/ref/spec#Function_declarations

			c.scanBlockStart()
			label := string(src[file.Offset(pos):file.Offset(c.Pos())])
			c.scanBlockEnd()
			to := c.At(file) + 1
			sections = append(sections, newSection(from, to, label))
		}
		if c.Token() != token.COMMENT {
			from = -1
		}
		lastTok = tok
	}
	// insert missing sections
	res := make([]Section, 0)
	if len(sections) == 0 {
		return res
	}
	first := sections[0]
	other := newOther(0, first.From(), src)
	if other.IsEmpty(src) {
		sections[0].from = 0
	} else {
		res = append(res, other)
	}

	for i := 0; i < len(sections)-1; i++ {
		a := sections[i]
		b := sections[i+1]
		c := newOther(a.To(), b.From(), src) // between
		if c.IsEmpty(src) {
			a.to = b.from
			res = append(res, a)
		} else {
			res = append(res, a, c)
		}
	}

	last := sections[len(sections)-1]
	res = append(res, last)

	if last.To() != len(src) {
		other := newOther(last.To(), len(src), src)
		if other.IsEmpty(src) {
			res[len(res)-1].to = len(src)
		} else {
			res = append(res, other)
		}
	}

	return res
}

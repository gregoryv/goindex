package goindex

import (
	"bytes"
	"fmt"
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

	var (
		c        = NewCursor(&s)
		sections = make([]Section, 0)
		from     int // sections start position within src
	)

	for c.Next() {
		pos := c.Pos()
		tok := c.Token()
		lit := c.Lit()

		if tok == token.COMMENT {
			if from != -1 {
				continue // multiline comment
			}
			// save position incase it's a related block comment
			from = file.Offset(pos)

			// single line comments come separately, check if there
			// are more
			if len(lit) > 2 && lit[:2] == "//" {
				continue
			}
		}

		// if comment blocks are followed by empty line, add them as
		// decoupled sections
		if tok != token.COMMENT && from != -1 {
			to := c.At(file) - 2
			if to > 0 && src[to] == '\n' {
				s := newOther(from, to, src)
				s.line = bytes.Count(src[:from], newLine) + 1
				sections = append(sections, s)
				from = -1
			}
		}

		if from == -1 { // no related comment
			from = file.Offset(pos)
		}
		// all cases here should scan the complete block and add a
		// section. Only interesting sections are added here, other
		// such as blocks of constants or variables are later filled
		// in as other sections.
		switch tok {
		case token.IMPORT:
			c.scanParenBlock()
			to := c.At(file) + 1
			s := newImport(from, to)
			s.line = bytes.Count(src[:file.Offset(pos)], newLine) + 1
			sections = append(sections, s)

		case token.TYPE:
			c.scanBlockStart()
			label := string(src[file.Offset(pos):file.Offset(c.Pos())])
			c.scanBlockEnd()
			to := c.At(file) + 1
			s := newSection(from, to, label)
			s.line = bytes.Count(src[:file.Offset(pos)], newLine) + 1
			sections = append(sections, s)

		case token.FUNC:
			// Fixme: func extra(), ie. no body
			// See https://go.dev/ref/spec#Function_declarations

			c.scanBlockStart()
			label := string(src[file.Offset(pos):file.Offset(c.Pos())])
			c.scanBlockEnd()
			to := c.At(file) + 1
			s := newSection(from, to, label)
			s.line = bytes.Count(src[:file.Offset(pos)], newLine) + 1
			sections = append(sections, s)
		}
		if c.Token() != token.COMMENT {
			from = -1
		}
	}
	// add any ending comments as their own block
	if from > -1 {
		sections = append(sections, newOther(from, len(src), src))
	}

	// insert missing sections, such as blocks of global const or var
	// declarations
	res := make([]Section, 0)
	if len(sections) == 0 {
		s := newOther(0, len(src), src)
		s.line = 1
		res = append(res, s)
		return res
	}
	first := sections[0]
	other := newOther(0, first.From(), src)
	if other.IsEmpty(src) {
		sections[0].from = 0
		sections[0].line = 1
	} else {
		other.line = 1
		res = append(res, other)
	}

	for i := 0; i < len(sections)-1; i++ {
		a := sections[i]
		b := sections[i+1]
		c := newOther(a.To(), b.From(), src) // between
		if c.IsEmpty(src) {
			a.to = b.from
			if a.line == 0 {
				fmt.Println("HERE")
			}
			res = append(res, a)
		} else {
			c.line = bytes.Count(src[:a.To()], newLine) + 1
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
			other.line = bytes.Count(src[:last.To()], newLine) + 1
			res = append(res, other)
		}
	}

	if i := len(res) - 1; res[i].line == 0 {
		res[i].line = bytes.Count(src[:res[i].From()], newLine) + 1
	}
	return res
}

var newLine = []byte("\n")

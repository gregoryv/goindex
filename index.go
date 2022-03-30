package goindex

import (
	"bytes"
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
	for c.Next() {
		pos := c.Pos()
		switch c.Token() {
		case token.COMMENT:
			from = file.Offset(pos) // and position to include in func blocks
			l := len(c.Lit())
			end := from + l + 5
			if end >= len(src) {
				end = len(src)
			}
			prefix := string(src[from+l+1 : end]) // if related it's either func or type
			//fmt.Printf("l=%v %q\n", l, prefix)
			if prefix != "func" {
				from = -1
			}

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

// ----------------------------------------

func newImport(from, to int) Section {
	return newSection(from, to, "import")
}

func newOther(from, to int, src []byte) Section {
	if to < from {
		to = from
	}
	part := bytes.TrimSpace(src[from:to])
	i := bytes.Index(part, []byte("\n"))
	var label string
	if i > -1 {
		label = string(part[:i])
	} else {
		label = string(part)
	}
	return newSection(from, to, label)
}

func newSection(from, to int, label string) Section {
	return Section{
		from:  from,
		to:    to,
		label: label,
	}
}

type Section struct {
	from, to int

	label string
}

func (me *Section) String() string         { return me.label }
func (me *Section) From() int              { return me.from }
func (me *Section) To() int                { return me.to }
func (me *Section) Grab(src []byte) []byte { return src[me.from:me.to] }
func (me *Section) IsEmpty(src []byte) bool {
	v := bytes.TrimSpace(me.Grab(src))
	return len(v) == 0
}

// ----------------------------------------

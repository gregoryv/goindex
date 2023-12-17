package goindex

import "bytes"

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

// Section defines a section within a Go source file
type Section struct {
	line, from, to int

	label string
}

// Line returns the line where this section starts. Lines begin with
// number 1.
func (me *Section) Line() int { return me.line }

// String returns short value of this section, e.g. for funcs only the
// signature
func (me *Section) String() string { return me.label }

// From returns the starting position of this section
func (me *Section) From() int { return me.from }

// To returns the end position of this section
func (me *Section) To() int { return me.to }

// Grab returns the the sections src[From:To]
func (me *Section) Grab(src []byte) []byte { return src[me.from:me.to] }

// IsEmpty returns true if the sections has no characters after
// bytes.TrimSpace has been applied.
func (me *Section) IsEmpty(src []byte) bool {
	v := bytes.TrimSpace(me.Grab(src))
	return len(v) == 0
}

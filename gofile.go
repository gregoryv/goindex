package gosort

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gregoryv/nexus"
)

func ParseGoFile(src []byte) *GoFile {
	return &GoFile{
		src:      src,
		sections: ParseSource(src),
	}
}

func ParseSource(src []byte) []Section {
	sections := make([]Section, 0)

	separators := []string{"func", "type", "//", "\n}"}
	for _, sep := range separators {
		index := indexAll(src, []byte(sep))
		for _, position := range index {
			sections = append(sections, Section{
				src:      src,
				position: position,
				ident:    strings.TrimSpace(sep),
			})
		}
	}
	sort.Sort(byPosition(sections))
	return sections
}

type GoFile struct {
	src      []byte
	sections []Section
}

func (me *GoFile) WriteTo(w io.Writer) (int64, error) {
	var last int
	p, err := nexus.NewPrinter(w)
	for _, s := range me.sections {
		p.Write(me.src[last:s.Position()])
		last = s.Position()
	}
	p.Write(me.src[last:]) // write final
	return p.Written, *err
}

func indexAll(src, sep []byte) []int {
	index := []int{}
	var last int
	for {
		i := bytes.Index(src[last:], sep)
		if i == -1 {
			break
		}
		index = append(index, last+i)
		last = last + i + len(sep)
	}
	return index
}

type byPosition []Section

func (s byPosition) Len() int      { return len(s) }
func (s byPosition) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byPosition) Less(i, j int) bool {
	return s[i].Position() < s[j].Position()
}

// ----------------------------------------

type Section struct {
	src      []byte
	position int
	ident    string
}

func (me *Section) Position() int { return me.position }

func (me *Section) String() string {
	return fmt.Sprintf("%d %s", me.position, me.ident)
}

func (me *Section) IsMethod() bool {
	if me.ident != "func" {
		return false
	}
	i := me.Position() + 5 // func (
	return me.src[i] == '('
}

func (me *Section) IsFunc() bool {
	if me.ident != "func" {
		return false
	}
	return true
}

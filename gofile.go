package gosort

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
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

func (me *GoFile) WriteTo(w io.Writer) (int, error) {
	var last int
	var total int
	for _, s := range me.sections {
		n, err := w.Write(me.src[last:s.Position()])
		if err != nil {
			return total + n, err
		}
		total += n
		last = s.Position()
	}
	n, err := w.Write(me.src[last:]) // write final
	return total + n, err
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
	position int
	ident    string
}

func (me *Section) Position() int { return me.position }

func (me *Section) String() string {
	return fmt.Sprintf("%d %s", me.position, me.ident)
}

package gosort

import (
	"bytes"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Test_indexFile(t *testing.T) {
	src, err := os.ReadFile("testdata/example.go")
	if err != nil {
		t.Fatal(err)
	}

	sections := ParseSource(src)
	var buf bytes.Buffer
	var last int
	for _, s := range sections {
		buf.Write(src[last:s.Position()])
		last = s.Position()
		t.Log(s.String())
	}
	buf.Write(src[last:]) // write final

	got, exp := buf.String(), string(src)

	if got != exp {
		golden.AssertEquals(t, got, exp)
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

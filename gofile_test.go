package gosort

import (
	"bytes"
	"os"
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

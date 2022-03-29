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

	var buf bytes.Buffer
	ParseGoFile(src).WriteTo(&buf)

	got, exp := buf.String(), string(src)
	if got != exp {
		golden.AssertEquals(t, got, exp)
	}
}

func TestSection_String(t *testing.T) {
	s := Section{
		position: 1,
		ident:    "x",
	}
	if got := s.String(); got == "" {
		t.Fail()
	}
}

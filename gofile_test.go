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
	gf := ParseGoFile(src)
	gf.WriteTo(&buf)

	got, exp := buf.String(), string(src)
	if got != exp {
		golden.AssertEquals(t, got, exp)
	}
}

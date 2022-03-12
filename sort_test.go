package gosort

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func Test_Index(t *testing.T) {
	src := []byte(`package x

// X stores info
type X struct {
}

// not much here

func play() {
_ = "hey"
}`)
	var from int
	var buf bytes.Buffer
	for _, to := range Index(src) {
		fmt.Fprint(&buf, string(src[from:to]))
		from = to
	}
	if got := buf.String(); got != string(src) {
		t.Error(got)
	}
}

func Example_FindConstructors() {
	src := []byte(`package x
type X struct {
}
func play() {}
func New() *X { return &X{} }
func Parse(v []byte) (*X, error) { return &X{}, nil }
func (c *X) Clone(_ interface{ Y() }) *X { t := *c; return &t}
`)

	fmt.Print(FindConstructors(src, "X"))
	// output:
	// [New Parse]
}

func Example_FindTypes() {
	src := []byte(`package x
type X struct {}
func NewX() *X { return &X{} }`)

	fmt.Print(FindTypes(src))
	// output:
	// [X]
}

func init() {
	log.SetFlags(0)
}

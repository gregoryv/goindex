package gosort

import (
	"fmt"
	"log"
)

func Example_FindConstructors() {
	src := []byte(`package x
type X struct {}
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

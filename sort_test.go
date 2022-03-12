package gosort

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(0)
}

func Example_FindTypes() {
	src := []byte(`package x

type X struct {}

func NewX() *X { return &X{} }
`)
	fmt.Print(FindTypes(src))
	// output:
	// [X]
}

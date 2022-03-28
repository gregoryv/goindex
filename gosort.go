package gosort

import (
	"fmt"
	"io"
)

func New(dst io.Writer, src []byte) *GoSort {
	return &GoSort{
		dst: dst,
		src: src,
	}
}

type GoSort struct {
	dst io.Writer
	src []byte
}

func (me *GoSort) Run() error {
	return fmt.Errorf("GoSort.Run: todo")
}

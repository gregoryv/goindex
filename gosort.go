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
	blocks := ParseBlocks(me.src)
	SortBlocks(blocks)

	for _, b := range blocks {
		b.WriteTo(me.dst)
		fmt.Fprintln(me.dst)
	}
	return nil
}

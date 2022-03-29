package gosort

import (
	"fmt"
	"io"
)

type GoFile struct {
	sections []Section
	src      []byte
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

// ----------------------------------------

type Section struct {
	position int
	ident    string
}

func (me *Section) Position() int { return me.position }

func (me *Section) String() string {
	return fmt.Sprintf("%d %s", me.position, me.ident)
}

package gosort

import "fmt"

type Section struct {
	position int
	ident    string
}

func (me *Section) Position() int { return me.position }

func (me *Section) String() string {
	return fmt.Sprintf("%d %s", me.position, me.ident)
}

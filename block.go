package gosort

import (
	"go/scanner"
	"go/token"
	"io"
)

func ParseBlocks(src []byte) []Block {
	result := make([]Block, 0)
	var from int
	for _, to := range Index(src) {
		result = append(result, Block{src[from:to]})
		from = to
	}
	return result
}

type Block struct {
	src []byte
}

func (me *Block) WriteTo(w io.Writer) {
	w.Write(me.src)
}

func (me *Block) IsConstructor(typeName string) bool {
	c := FindConstructors(me.src, typeName)
	return len(c) > 0
}

func (me *Block) IsMethod(typeName string) bool {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(me.src))
	s.Init(file, me.src, nil /* no error handler */, 0)

loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break loop
		case token.FUNC:
			_, tok, _ := s.Scan()
			if tok == token.LPAREN { // ie. method
				return true
			}
		}
	}
	return false
}

func (me *Block) IsType(typeName string) bool {
	for _, t := range FindTypes(me.src) {
		if t == typeName {
			return true
		}
	}
	return false
}

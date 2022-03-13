package gosort

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io"
)

func ParseBlocks(src []byte) []Block {
	result := make([]Block, 0)
	var from int
	for _, to := range Index(src) {
		result = append(result, NewBlock(src[from:to]))
		from = to
	}
	return result
}

func NewBlock(src []byte) Block {
	b := Block{
		src: src,
	}
	// parse name of type, constructor or method
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, 0)

loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break loop

		case token.TYPE:
			_, _, lit := s.Scan()
			b.name = lit
			b.rel = lit
			b.isType = true // mark

		case token.FUNC:
			_, tok, lit := s.Scan()
			switch tok {

			case token.LPAREN: // ie. method
				b.isMethod = true // mark
				b.rel = receiverName(&s)
				_, _, lit := s.Scan()
				b.name = lit // name of method

			case token.IDENT: // func or constructor
				b.name = lit
				skipFuncArgs(&s)
				b.rel = returnName(&s)
			}
		}
	}
	return b
}

func receiverName(s *scanner.Scanner) string {
	var name string
loop:
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			break loop

		case token.IDENT:
			name = lit

		case token.RPAREN:
			break loop
		}
	}
	return name
}

func skipFuncArgs(s *scanner.Scanner) {
	var left int
loop:
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break loop

		case token.LPAREN:
			left++

		case token.RPAREN:
			left--
			if left == 0 {
				return
			}
		}
	}
}

func returnName(s *scanner.Scanner) string {
	var name string
loop:
	for {
		_, tok, lit := s.Scan()
		switch tok {
		case token.EOF:
			break loop

		case token.IDENT:
			name = lit

		case token.LBRACE:
			break loop
		}
	}
	return name
}

type Block struct {
	src []byte

	name string // if type, method or constructor

	// rel is the type name this block is related to
	// for methods and constructors it's the name of the type
	// for type blocks it's the same as the name field
	rel string

	isType bool

	isMethod      bool
	isConstructor bool
	isFunc        bool
}

func (me *Block) String() string {
	return fmt.Sprintf("%s %s", me.rel, me.name)
}

func (me *Block) WriteTo(w io.Writer) {
	w.Write(me.src)
}

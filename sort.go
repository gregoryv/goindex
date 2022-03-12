package gosort

import (
	"go/scanner"
	"go/token"
)

// FindTypes returns all named types in the given Go source
func FindTypes(src []byte) []string {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, 0)

	result := make([]string, 0)
	for {
		_, tok, _ := s.Scan()
		switch tok {
		case token.EOF:
			break
		case token.TYPE:
			_, _, name := s.Scan()
			result = append(result, name)
		}
	}
	return result
}

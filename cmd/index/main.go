// Command index indexes source files into sections
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gregoryv/goindex"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s FILES...\n", os.Args[0])
		flag.PrintDefaults()
	}
	verbose := flag.Bool("verbose", false, "progress to stderr")
	flag.Parse()

	for _, filename := range flag.Args() {
		ext := filepath.Ext(filename)
		switch ext  {
		case ".go":
			if *verbose {
				fmt.Fprintln(os.Stderr, "parse", filename, ext)
			}
			src, err := os.ReadFile(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			sections := goindex.Index(src)
			for _, s := range sections {
				fmt.Printf("%s %s\n", filename, s.String())
			}
		default:
			if *verbose {
				fmt.Fprintln(os.Stderr, "skip", filename, ext)
			}
		}
	}
}

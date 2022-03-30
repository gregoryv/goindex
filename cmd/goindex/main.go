package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gregoryv/goindex"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	for _, filename := range flag.Args() {
		src, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		sections := goindex.Index(src)
		for _, s := range sections {
			fmt.Printf("%s %v %v %s\n", filename, s.From(), s.To(), s.String())
		}
	}
}

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
	filename := flag.Arg(0)

	log.SetFlags(0)

	src, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	sections := goindex.Index(src)
	for _, s := range sections {
		fmt.Printf("%v %v %s\n", s.From(), s.To(), s.String())
	}
}

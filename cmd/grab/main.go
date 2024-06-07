// Command grab extracts sections from a file by byte range
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		w := os.Stderr
		fmt.Fprintf(w, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprint(w, `
grab command reads frm stdin in the form

FILE1 FROM TO
FILE1 FROM TO
FILE2 FROM TO

and writes out the sections from those files to stdout.

FROM and TO are the byte index in each file.

Options

`)
		flag.PrintDefaults()
	}

	cut := flag.Bool("cut", false, "")
	flag.BoolVar(cut, "c", *cut, "cut grabbed sections from source file")

	flag.Parse()

	log.SetFlags(0)

	s := bufio.NewScanner(os.Stdin)
	g := Grabber{cut: *cut}

	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) < 3 {
			log.Println(line)
			log.Fatal("line should be: FILE FROM TO")
		}
		if err := g.Grab(os.Stdout, f[0], f[1], f[2]); err != nil {
			log.Fatal(err)
		}
	}
	if err := g.Flush(); err != nil {
		log.Fatal(err)
	}
}

type Grabber struct {
	filename string
	src      []byte

	cut      bool
	modified []byte

	lastTo int
}

func (g *Grabber) Grab(w io.Writer, filename, sfrom, sto string) error {
	// Load src from new file
	if g.filename != filename {

		if g.filename != "" && g.cut {
			if err := g.Flush(); err != nil {
				return err
			}
		}
		src, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		g.src = src
		g.filename = filename
	}
	from, err := strconv.Atoi(sfrom)
	if err != nil {
		return err
	}
	to, err := strconv.Atoi(sto)
	if err != nil {
		return err
	}
	w.Write(g.src[from:to])
	// save rest for later write
	if g.cut {
		g.modified = append(g.modified, g.src[g.lastTo:from]...)
		g.lastTo = to
	}
	return nil
}

func (g *Grabber) Flush() error {
	if !g.cut || g.filename == "" {
		return nil
	}
	// add tail data
	g.modified = append(g.modified, g.src[g.lastTo:]...)
	fi, _ := os.Stat(g.filename)
	if err := os.WriteFile(g.filename, g.modified, fi.Mode()); err != nil {
		return err
	}
	// reset
	g.modified = g.modified[0:0]
	g.lastTo = 0
	return nil
}

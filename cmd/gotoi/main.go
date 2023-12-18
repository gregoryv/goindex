// Command goto opens file on specific line
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gregoryv/vt100"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gotoi [OPTIONS] INDEX")
		flag.PrintDefaults()
	}
	filename := flag.String("f", ".index", "goindex file")
	flag.Parse()

	args := flag.Args()

	index := make(map[int]bool)
	for _, a := range args {
		i, err := strconv.Atoi(a)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		index[i] = true
	}
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	in := bufio.NewScanner(bytes.NewReader(data))
	var i int
	var lastFile string
	for in.Scan() {
		i++
		line := in.Text()
		part := strings.SplitN(line, " ", 5)
		if len(part) < 5 {
			fmt.Fprintln(os.Stderr, line)
			fmt.Fprintln(os.Stderr, "line should be: FILE FROM TO LINE ...")
			continue
		}
		switch {
		case len(args) == 0:
			if part[0] != lastFile {
				lastFile = part[0]
				fmt.Print(fg.Cyan.String()+attr.Dim.String()+strings.TrimSpace(lastFile), attr.Reset, "\n")
			}
			fmt.Printf("%v %s\n", i, paint(part[4]))

		case index[i]:
			exec.Command("emacsclient", "-n", "+"+part[3], part[0]).Run()
		}
	}
}

func paint(v string) string {
	var sb strings.Builder
	first := v
	i := strings.Index(v, " ")
	if i > 0 {
		first = first[:i]
	} else {
		i = 0
	}
	switch first {
	case "import":
		sb.WriteString(fg.Magenta.String())
		sb.WriteString(first)
		sb.WriteString(attr.Reset.String())
		return sb.String()

	case "func":
		sb.WriteString(fg.Magenta.String())
		sb.WriteString(first)
		sb.WriteString(attr.Reset.String())
		sb.WriteString(v[i:])
		return sb.String()

	case "var", "const", "type", "package":
		sb.WriteString(fg.Magenta.String())
		sb.WriteString(first)
		sb.WriteString(attr.Reset.String())
		sb.WriteString(v[i:])
		return sb.String()

	case "//":
		sb.WriteString(fg.Green.String())
		sb.WriteString(attr.Dim.String())
		sb.WriteString(v)
		sb.WriteString(attr.Reset.String())
		return sb.String()

	default:
		return v
	}
}

var (
	fg   = vt100.ForegroundColors()
	attr = vt100.Attributes()
)

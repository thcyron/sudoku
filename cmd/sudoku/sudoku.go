package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/thcyron/sudoku"
)

func main() {
	var f *os.File

	switch len(os.Args) {
	case 1:
		f = os.Stdin
	case 2:
		ff, err := os.Open(os.Args[1])
		if err != nil {
			die("%v", err)
		}
		f = ff
	default:
		fmt.Fprintf(os.Stderr, "usage: sudoku [file]\n")
		os.Exit(2)
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		die("%v", err)
	}

	var s string
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if strings.IndexByte("0123456789_ ", b) >= 0 {
			s += string(b)
		}
	}
	if len(s) != 9*9 {
		die("incomplete grid")
	}

	grid := sudoku.NewGrid()
	if ok, err := grid.Set(s); !ok {
		die("%v", err)
	}
	if grid.Solve() {
		grid.Print()
	} else {
		die("not solvable")
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("sudoku: %s\n", format))
	os.Exit(1)
}

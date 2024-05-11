package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/f01c33/strsearch"
)

func main() {
	e := ""
	f := ""
	// g := ""
	p := ""
	flag.StringVar(&e, "E", "", "go regex")
	// flag.StringVar(&g, "g", "", "glob match")
	flag.StringVar(&p, "p", "", "pattern to find (regular text, no regex or similar)")
	flag.StringVar(&f, "f", "", "file to read, defaults to stdin")
	flag.Parse()
	var file io.Reader = os.Stdin
	if f != "" {
		fl, err := os.Open(f)
		if err != nil {
			fmt.Println("File not found")
		}
		defer fl.Close()
		file = fl
	}
	rdr := bufio.NewReader(file)

	rx := regexp.MustCompile(e)
	var err error
	var line []byte
	var idxs [][]int
	for {
		line, _, err = rdr.ReadLine()
		if err != nil {
			break
		}
		if e != "" {
			idxs = rx.FindAllIndex(line, -1)
		} else if p != "" {
			idxs = strsearch.FindAllIndex(line, []byte(p))
		}
		if idxs == nil {
			continue
		}
		curr1 := 0
		i := 0
		Reset := "\033[0m"
		Red := "\033[31m"
		for {
			if curr1 == len(idxs) {
				fmt.Printf("%s\n", string(line[idxs[curr1-1][1]:]))
				break
			}
			if i < idxs[curr1][0] {
				fmt.Print(string(line[i:idxs[curr1][0]]))
				fmt.Print(Red, string(line[idxs[curr1][0]:idxs[curr1][1]]), Reset)
				i = idxs[curr1][1]
				curr1++
				continue
			}
		}
	}
}

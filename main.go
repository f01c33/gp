package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/f01c33/strsearch"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func main() {
	e := ""
	f := ""
	// g := ""
	p := ""
	fz := ""
	flag.StringVar(&e, "E", "", "go regex")
	// flag.StringVar(&g, "g", "", "glob match")
	flag.StringVar(&p, "p", "", "pattern to find (regular text, no regex or similar)")
	flag.StringVar(&fz, "fz", "", "fuzzy search")
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
		} else if fz != "" {
			strList := strings.Split(string(line), " ")
			ranks := fuzzy.RankFindNormalized(fz, strList)
			if len(ranks) > 0 {
				idxs = [][]int{}
			} else {
				idxs = nil
			}
			ri := 0
			currRank := 0
			for _, s := range strList {
				ri += len(s) + 1
				if currRank == len(ranks) {
					break
				}
				if ranks[currRank].Target == s {
					idxs = append(idxs, []int{ri - len(s) - 1, ri - 1})
					currRank++
				}
			}
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

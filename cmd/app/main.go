package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/clsrfish/invest/internal/download"
	"github.com/clsrfish/invest/internal/parser"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

const (
	defaultInputFile  = "./input.txt"
	defaultOutputFile = "./books.csv"
)

type SortByVal []lo.Entry[string, int]

func (a SortByVal) Len() int      { return len(a) }
func (a SortByVal) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByVal) Less(i, j int) bool {
	if a[i].Value != a[j].Value {
		return a[i].Value < a[j].Value
	}
	return a[i].Key < a[j].Key
}

var (
	inputf string
	outputf string
)

func main() {

	flag.StringVar(&inputf, "i", defaultInputFile, "Input file")
	flag.StringVar(&outputf, "o", defaultOutputFile, "Output file")
	flag.Parse()

	sites := make(map[string]int)
	if input, err := os.Open(inputf); err != nil {
		panic(err)
	} else {
		sc := bufio.NewScanner(input)
		for sc.Scan() {
			line := sc.Text()
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//") {
				continue
			}

			if _, err := url.Parse(line); err != nil {
				log.Warnf("invalid url: %s\n", line)
				continue
			}
			sites[line] = 1
		}
	}

	recommendation := make(map[string]int, 0)

	lo.ForEach(lo.Keys(sites), func(site string, i int) {
		f, err := download.Download(site)
		if err != nil {
			panic(err)
		}
		if bs, err := ioutil.ReadFile(f); err != nil {
			panic(err)
		} else {
			books := parser.ParseBooks(string(bs))
			for book := range books {
				recommendation[book] += 1
			}
		}
	})

	res := lo.Entries(recommendation)
	sort.Sort(SortByVal(res))
	res = lo.Reverse(res)
	log.Infof("total: %d\n", len(res))
	if output, err := os.OpenFile(outputf, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		panic(err)
	} else {
		output.WriteString("book, recommends\n")
		for _, entry := range res {
			output.WriteString(fmt.Sprintf("\"%s\", %d\n", entry.Key, entry.Value))
		}
	}
}

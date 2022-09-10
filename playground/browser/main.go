package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/browser"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [url]\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := flag.Args()
	switch len(args) {
	case 0:
		check(browser.OpenReader(os.Stdin))
	case 1:
		check(browser.OpenURL(args[0]))
	default:
		usage()
	}
}

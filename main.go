// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dex/iso9660"
)

var port = flag.Int("p", 8080, "`Port` to linsten on")

var flag.Usage = func () {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <iso>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

type isoFs struct {
	fs *iso9660.FileSystem
}

func (iso *isoFs) Open(name string) (http.File, error) {
	return iso.fs.Open(name)
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fs, err := iso9660.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), http.FileServer(&isoFs{fs})))
}

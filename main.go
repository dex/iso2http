// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dex/iso9660"
)

var port = flag.Int("p", 8080, "`Port` to linsten on")

var Usage = func() {
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

func download(url string) (fileName string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create temp file
	tmpfile, err := ioutil.TempFile("", "*.iso")
	if err != nil {
		return
	}

	// Write the body to file
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return
	}
	fileName = tmpfile.Name()
	return
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	isofile := flag.Arg(0)

	if strings.HasPrefix(isofile, "http") {
		fmt.Println("Download ISO file")
		tmpName, err := download(isofile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Download completed")
		isofile = tmpName
		defer os.Remove(tmpName)
	}

	fs, err := iso9660.Open(isofile)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), http.FileServer(&isoFs{fs})))
}

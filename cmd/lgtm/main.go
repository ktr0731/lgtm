package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ktr0731/lgtm"
)

var (
	w = flag.Bool("w", false, "write to input file")
)

func main() {
	flag.Parse()

	switch {
	case *w && len(flag.Args()) != 1:
		fmt.Println("usage: lgtm -w <input.gif>")
		os.Exit(1)
	case !*w && len(flag.Args()) != 2:
		fmt.Println("usage: lgtm <input.gif> <output.gif>")
		os.Exit(1)
	}

	var r io.Reader
	var err error
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r = f

	var outFileName string
	if *w {
		buf := new(bytes.Buffer)
		io.Copy(buf, f)
		r = buf
		f.Close()

		outFileName = flag.Arg(0)
	} else {
		outFileName = flag.Arg(1)
	}

	w, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err := lgtm.New(r, w); err != nil {
		panic(err)
	}
}

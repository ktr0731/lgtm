package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ktr0731/lgtm"
	"golang.org/x/sync/errgroup"
)

var (
	w = flag.Bool("w", false, "write to input file")
	t = flag.Float64("threshold", 0.3, "margin threshold")
)

func main() {
	flag.Parse()

	lgtm.Threshold = *t

	switch {
	case *w && len(flag.Args()) > 0:
		var eg errgroup.Group
		for _, a := range flag.Args() {
			a := a
			eg.Go(func() error {
				return overlayLGTM(a, a)
			})
		}
		if err := eg.Wait(); err != nil {
			log.Fatal(err)
		}
	case *w && len(flag.Args()) == 0:
		fmt.Println("usage: lgtm -w <input.gif>")
		os.Exit(1)
	case !*w && len(flag.Args()) != 2:
		fmt.Println("usage: lgtm <input.gif> <output.gif>")
		os.Exit(1)
	default:
		if err := overlayLGTM(flag.Arg(1), flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}
}

func overlayLGTM(outFile, inFile string) error {
	var err error
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}

	r := new(bytes.Buffer)
	io.Copy(r, f)
	f.Close()

	w, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer w.Close()

	return lgtm.New(r, w)
}

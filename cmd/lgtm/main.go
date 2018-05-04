package main

import (
	"os"

	"github.com/ktr0731/lgtm"
)

func main() {
	r, err := os.Open("example.gif")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	w, err := os.Create("out.gif")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err := lgtm.New(r, w); err != nil {
		panic(err)
	}
}

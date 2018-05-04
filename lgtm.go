package lgtm

import (
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

var lgtm image.Image
var once sync.Once

func initBaseImage() error {
	var err error
	once.Do(func() {
		fs, err := fs.New()
		if err != nil {
			err = errors.Wrap(err, "failed to setup LGTM image")
			return
		}
		f, err := fs.Open("lgtm.png")
		if err != nil {
			err = errors.Wrap(err, "failed to open lgtm.png")
			return
		}
		defer f.Close()
		lgtm, err = png.Decode(f)
		if err != nil {
			err = errors.Wrap(err, "failed to decode lgtm.png to PNG image")
			return
		}
	})
	return err
}

func NewLGTM(r io.Reader, w io.Writer) error {
	if err := initBaseImage(); err != nil {
		return errors.Wrap(err, "failed to init base image")
	}

	g, err := gif.DecodeAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode source GIF")
	}

	// outGIF := &gif.GIF{
	// 	Image: make([]*image.Paletted, 0, len(inGIF.Image)),
	// 	Delay: inGIF.Delay,
	// }
	for i := range g.Image {
		// TODO: not ZP
		draw.Draw(g.Image[i], g.Image[i].Bounds(), lgtm, image.ZP, draw.Over)
	}

	return gif.EncodeAll(w, g)
}

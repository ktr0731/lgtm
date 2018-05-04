package lgtm

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"sync"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"

	_ "github.com/ktr0731/lgtm/statik"
	"github.com/rakyll/statik/fs"
)

var lgtm image.Image
var once sync.Once

func initBaseImage() error {
	var werr error
	once.Do(func() {
		sfs, err := fs.New()
		if err != nil {
			werr = errors.Wrap(err, "failed to setup LGTM image")
			return
		}
		f, err := sfs.Open("/lgtm.png")
		if err != nil {
			werr = errors.Wrap(err, "failed to open lgtm.png")
			return
		}
		defer f.Close()
		lgtm, err = png.Decode(f)
		if err != nil {
			werr = errors.Wrap(err, "failed to decode lgtm.png to PNG image")
			return
		}
	})
	return werr
}

func New(r io.Reader, w io.Writer) error {
	if err := initBaseImage(); err != nil {
		return errors.Wrap(err, "failed to init base image")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to read from input")
	}

	if g, err := gif.DecodeAll(bytes.NewBuffer(b)); err == nil {
		return fromGIF(w, g)
	} else if p, err := png.Decode(bytes.NewBuffer(b)); err == nil {
		return fromImage(w, p)
	} else {
		return errors.Wrap(err, "unsupported file type")
	}
}

func fromGIF(w io.Writer, g *gif.GIF) error {
	lgtm := adjustedLGTM(g.Image[0].Bounds())
	for i := range g.Image {
		drawOver(g.Image[i], lgtm)
	}
	return gif.EncodeAll(w, g)
}

// only png
func fromImage(w io.Writer, img image.Image) error {
	lgtm := adjustedLGTM(img.Bounds())
	p := image.NewRGBA(img.Bounds())
	draw.Draw(p, img.Bounds(), img, image.ZP, draw.Src)
	drawOver(p, lgtm)
	return png.Encode(w, p)
}

func drawOver(img draw.Image, lgtm image.Image) {
	p := lgtm.Bounds().Size()
	p.X = -((img.Bounds().Dx() - p.X) / 2)
	p.Y = -((img.Bounds().Dy() - p.Y) / 2)
	draw.Draw(img, img.Bounds(), lgtm, p, draw.Over)
}

func adjustedLGTM(r image.Rectangle) image.Image {
	b := lgtm.Bounds()
	threshold := 0.3

	var x, y uint
	if b.Dx() != r.Dx() {
		ratio := float64(r.Dx()) / float64(b.Dx())
		fx := math.Floor(float64(b.Dx()) * ratio)
		fy := math.Floor(float64(b.Dy()) * ratio)
		x = uint(fx - fx*threshold)
		y = uint(fy - fy*threshold)
	} else if b.Dy() > r.Dy() {
		ratio := float64(r.Dy()) / float64(b.Dy())
		fx := math.Floor(float64(b.Dx()) * ratio)
		fy := math.Floor(float64(b.Dy()) * ratio)
		x = uint(fx - fx*threshold)
		y = uint(fy - fy*threshold)
	}
	return resize.Resize(x, y, lgtm, resize.Lanczos3)
}

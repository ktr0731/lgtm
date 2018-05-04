package lgtm

import (
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
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

	g, err := gif.DecodeAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode source GIF")
	}

	// outGIF := &gif.GIF{
	// 	Image: make([]*image.Paletted, 0, len(inGIF.Image)),
	// 	Delay: inGIF.Delay,
	// }
	lgtm := adjustedLGTM(g.Image[0])
	bounds := g.Image[0].Bounds()
	for i := range g.Image {
		p := lgtm.Bounds().Size()
		p.X = -((bounds.Dx() - p.X) / 2)
		p.Y = -((bounds.Dy() - p.Y) / 2)
		draw.Draw(g.Image[i], bounds, lgtm, p, draw.Over)
	}

	return gif.EncodeAll(w, g)
}

func adjustedLGTM(p *image.Paletted) image.Image {
	b := lgtm.Bounds()
	if b.Dx() <= p.Bounds().Dx() && b.Dy() <= p.Bounds().Dy() {
		return lgtm
	}

	threshold := 0.3

	var x, y uint
	if b.Dx() > p.Bounds().Dx() {
		ratio := float64(p.Bounds().Dx()) / float64(b.Dx())
		fx := math.Floor(float64(b.Dx()) * ratio)
		fy := math.Floor(float64(b.Dy()) * ratio)
		x = uint(fx - fx*threshold)
		y = uint(fy - fy*threshold)
	} else if b.Dy() > p.Bounds().Dy() {
		ratio := float64(p.Bounds().Dy()) / float64(b.Dy())
		fx := math.Floor(float64(b.Dx()) * ratio)
		fy := math.Floor(float64(b.Dy()) * ratio)
		x = uint(fx - fx*threshold)
		y = uint(fy - fy*threshold)
	}
	return resize.Resize(x, y, lgtm, resize.Lanczos3)
}

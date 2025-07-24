package service

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

type palette struct {
	bg, fg []string
}

func (p *palette) AddBackgroundColor(c string) { p.bg = append(p.bg, c) }
func (p *palette) AddForegroundColor(c string) { p.fg = append(p.fg, c) }

func (p palette) Sprint(output *termenv.Output) string {
	var b strings.Builder

	for i, c := range p.bg {
		fmt.Fprint(&b,
			output.String(" 1 ").
				Background(output.Color(c)).
				Foreground(output.Color(p.fg[i])),
			" ")
	}
	return b.String()
}

func (p palette) Index(i int) (bgColor, fgColor string) {
	if i >= 1 {
		i -= 1
	}
	return p.bg[i], p.fg[i]
}

// New palette of colors from startColor to endColor with number of shades (gradients).
func NewPaletter(startColor, endColor string, shades int) Paletter {
	c1, err := colorful.Hex(startColor)
	if err != nil {
		panic(err)
	}

	c2, err := colorful.Hex(endColor)
	if err != nil {
		panic(err)
	}

	adaptiveForeground := func(bg colorful.Color) colorful.Color {
		bh, _, bl := bg.Hcl()
		if bl > 0.2 {
			return colorful.Hcl(360-bh, 0.2, 1)
		}
		return colorful.Hcl(0, 0, 1)
	}

	var p palette
	for i := range shades {
		bg := c1.BlendRgb(c2, float64(i)/float64(shades-1))
		p.AddBackgroundColor(bg.Hex())
		fg := adaptiveForeground(colorful.LinearRgb(174, 136, 227))
		p.AddForegroundColor(fg.Hex())
	}

	return &p
}

package service

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

type palette []string

func (p palette) Sprint(output *termenv.Output) string {
	var b strings.Builder

	for _, c := range p {
		fmt.Fprint(&b, output.String("  ").Background(output.Color(c)), " ")
	}
	return b.String()
}

func (p palette) Index(i int) string {
	if i >= 1 {
		return p[i-1]
	}
	return p[i]
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

	var p palette
	for i := range shades {
		p = append(p, c1.BlendRgb(c2, float64(i)/float64(shades-1)).Hex())
	}

	return p
}

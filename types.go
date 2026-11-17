package maprenderer

import (
	"image/color"

	"github.com/luanti-golang/types"
)

type NodeWithColor struct {
	*types.Node
	Color *color.RGBA
}

// toNRGBA reinterprets the (non-premultiplied) color values of the color resolver as color.NRGBA
func toNRGBA(c *color.RGBA) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

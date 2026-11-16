package maprenderer

import (
	"image/color"

	"github.com/luanti-golang/types"
)

type NodeWithColor struct {
	*types.Node
	Color *color.RGBA
}

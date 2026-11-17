package maprenderer

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/luanti-golang/types"
)

func GetIsometricImageSize(size *types.Pos, cube_len int) (int, int) {
	width := (size.X() * cube_len / 2) +
		(size.Z() * cube_len / 2)

	height := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4)

	return width, height
}

func GetIsoCenterCubeOffset(size *types.Pos, cube_len int) (int, int) {
	x := (size.Z() * cube_len / 2) -
		(cube_len / 2)

	y := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4) -
		cube_len

	return x, y
}

func GetIsoCubePosition(center_x, center_y, cube_len int, pos *types.Pos) (int, int) {
	x := center_x -
		(pos.Z() * cube_len / 2) +
		(pos.X() * cube_len / 2)

	y := center_y -
		(pos.X() * cube_len / 4) -
		(pos.Y() * cube_len / 2) -
		(pos.Z() * cube_len / 4)

	return x, y
}

// premultipliedColor holds alpha-premultiplied 16-bit color components, as returned by color.Color.RGBA()
type premultipliedColor struct {
	r, g, b, a uint32
}

func newPremultipliedColor(c color.Color) premultipliedColor {
	r, g, b, a := c.RGBA()
	return premultipliedColor{r: r, g: g, b: b, a: a}
}

// premultiplyNRGBA converts a non-premultiplied color without going through the color.Color interface
func premultiplyNRGBA(c color.NRGBA) premultipliedColor {
	// same as color.NRGBA.RGBA()
	a := uint32(c.A) * 0x101
	return premultipliedColor{
		r: uint32(c.R) * 0x101 * a / 0xffff,
		g: uint32(c.G) * 0x101 * a / 0xffff,
		b: uint32(c.B) * 0x101 * a / 0xffff,
		a: a,
	}
}

// blendPixel composites the color over the existing pixel (porter-duff "over")
func blendPixel(img *image.RGBA, x, y int, c premultipliedColor) {
	if !(image.Point{X: x, Y: y}.In(img.Rect)) {
		return
	}
	i := img.PixOffset(x, y)
	pix := img.Pix[i : i+4 : i+4]

	if c.a == 0xffff {
		// opaque, plain overwrite
		pix[0] = uint8(c.r >> 8)
		pix[1] = uint8(c.g >> 8)
		pix[2] = uint8(c.b >> 8)
		pix[3] = uint8(c.a >> 8)
		return
	}

	// same formula as image/draw with draw.Over
	ia := 0xffff - c.a
	pix[0] = uint8((uint32(pix[0])*0x101*ia/0xffff + c.r) >> 8)
	pix[1] = uint8((uint32(pix[1])*0x101*ia/0xffff + c.g) >> 8)
	pix[2] = uint8((uint32(pix[2])*0x101*ia/0xffff + c.b) >> 8)
	pix[3] = uint8((uint32(pix[3])*0x101*ia/0xffff + c.a) >> 8)
}

// DrawIsoCube draws an isometric cube with its top-left corner at the given offset.
// The cube is composited over the existing image content ("over" operator), so translucent colors
// let the already drawn content shine through. Colors are interpreted according to their color.Color
// implementation, use color.NRGBA for non-premultiplied colors.
func DrawIsoCube(img *image.RGBA, cube_len, x_offset, y_offset int, c1, c2, c3 color.Color) error {
	err := validateCubeLen(cube_len)
	if err != nil {
		return err
	}
	drawIsoCube(img, cube_len, x_offset, y_offset, newPremultipliedColor(c1), newPremultipliedColor(c2), newPremultipliedColor(c3))
	return nil
}

func validateCubeLen(cube_len int) error {
	if cube_len%4 != 0 {
		return fmt.Errorf("cube_len must be divisible by 4")
	}
	if cube_len <= 4 {
		return fmt.Errorf("cube_len must be greater than 4")
	}
	return nil
}

// drawIsoCube draws the cube without validating the cube length
func drawIsoCube(img *image.RGBA, cube_len, x_offset, y_offset int, pc1, pc2, pc3 premultipliedColor) {
	half_len_zero_indexed := (cube_len / 2) - 1
	quarter_len := cube_len / 4

	// left/right part
	yo := 0
	for x := 0; x <= half_len_zero_indexed; x++ {
		for y := 0; y <= half_len_zero_indexed; y++ {
			// left
			blendPixel(img, x_offset+x, y_offset+y+quarter_len+yo, pc1)
			// right
			blendPixel(img, x_offset+cube_len-1-x, y_offset+y+quarter_len+yo, pc2)
		}
		if x%2 == 0 {
			yo = yo + 1
		}
	}

	// upper part
	yo = 0
	yl := 1
	for x := 0; x <= half_len_zero_indexed-1; x++ {
		for y := 0; y <= yl; y++ {
			// left
			blendPixel(img, x_offset+1+x, y_offset+quarter_len-1-yo+y, pc3)
			// right
			blendPixel(img, x_offset+cube_len-2-x, y_offset+quarter_len-1-yo+y, pc3)
		}
		if x%2 != 0 {
			yo = yo + 1
			yl = yl + 2
		}
	}

}

func GetIsoNodeOrder(p *types.Pos) int {
	return (64000 - p.X()) + p.Y() + (64000 - p.Z())
}

func addAndClampUint8(a uint8, b int) uint8 {
	v := int(a) + b
	if v > 255 {
		return 255
	} else if v < 0 {
		return 0
	} else {
		return uint8(v)
	}
}

// adjustNRGBA brightens or darkens the non-premultiplied color by the given value
func adjustNRGBA(c color.NRGBA, value int) color.NRGBA {
	return color.NRGBA{
		R: addAndClampUint8(c.R, value),
		G: addAndClampUint8(c.G, value),
		B: addAndClampUint8(c.B, value),
		A: c.A,
	}
}

func ColorAdjust(c *color.RGBA, value int) *color.RGBA {
	return &color.RGBA{
		R: addAndClampUint8(c.R, value),
		G: addAndClampUint8(c.G, value),
		B: addAndClampUint8(c.B, value),
		A: c.A,
	}
}

// BlendColor composites the non-premultiplied fg color over bg ("over" operator).
// The alpha of the foreground color is divided by bf to make stacked translucent nodes less dominant.
func BlendColor(bg, fg *color.RGBA, bf float64) *color.RGBA {
	fa := float64(fg.A) / 255 / bf
	ba := float64(bg.A) / 255
	// contribution of the background after being covered by the foreground
	bw := ba * (1 - fa)
	oa := fa + bw

	if oa <= 0 {
		return &color.RGBA{}
	}

	blend := func(f, b uint8) uint8 {
		return uint8((float64(f)*fa + float64(b)*bw) / oa)
	}

	return &color.RGBA{
		R: blend(fg.R, bg.R),
		G: blend(fg.G, bg.G),
		B: blend(fg.B, bg.B),
		A: uint8(oa*255 + 0.5),
	}
}

func SortNodesWithColor(n1, n2 *NodeWithColor) int {
	o1 := GetIsoNodeOrder(n1.Pos)
	o2 := GetIsoNodeOrder(n2.Pos)
	return o1 - o2
}

// sortNodesByIsoOrder sorts the nodes back-to-front for drawing (stable counting sort on GetIsoNodeOrder).
// The order values of a render span a small range (sum of the dimensions), so this is O(n + range).
func sortNodesByIsoOrder(nodes []*NodeWithColor) []*NodeWithColor {
	if len(nodes) < 2 {
		return nodes
	}

	orders := make([]int, len(nodes))
	min_order, max_order := math.MaxInt, math.MinInt
	for i, n := range nodes {
		o := GetIsoNodeOrder(n.Pos)
		orders[i] = o
		min_order = min(min_order, o)
		max_order = max(max_order, o)
	}

	// offsets[o] = start index of the order value o in the sorted output
	offsets := make([]int, max_order-min_order+2)
	for _, o := range orders {
		offsets[o-min_order+1]++
	}
	for i := 1; i < len(offsets); i++ {
		offsets[i] += offsets[i-1]
	}

	sorted := make([]*NodeWithColor, len(nodes))
	for i, n := range nodes {
		idx := orders[i] - min_order
		sorted[offsets[idx]] = n
		offsets[idx]++
	}
	return sorted
}

package maprenderer_test

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/luanti-golang/maprenderer"
	"github.com/luanti-golang/types"
	"github.com/stretchr/testify/assert"
)

// stone layer at y=0, water layer at y=1, nothing else
func layeredNodeAccessor(p *types.Pos) (*types.Node, error) {
	switch p.Y() {
	case 0:
		return &types.Node{Name: "stone", Pos: p}, nil
	case 1:
		return &types.Node{Name: "water", Pos: p}, nil
	}
	return nil, nil
}

func layeredColorResolver(name string, param2 int) *color.RGBA {
	switch name {
	case "stone":
		return &color.RGBA{R: 100, G: 100, B: 100, A: 255}
	case "water":
		return &color.RGBA{R: 0, G: 0, B: 200, A: 100}
	}
	return nil
}

func TestRenderIsometricSwappedCorners(t *testing.T) {
	p1 := types.NewPos(0, 0, 0)
	p2 := types.NewPos(7, 3, 7)

	expected, err := maprenderer.RenderIsometric(layeredNodeAccessor, layeredColorResolver, p1, p2, nil)
	assert.NoError(t, err)

	img, err := maprenderer.RenderIsometric(layeredNodeAccessor, layeredColorResolver, p2, p1, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, img)
}

func TestRenderIsometricTileSwappedCorners(t *testing.T) {
	p1 := types.NewPos(0, 0, 0)
	p2 := types.NewPos(7, 3, 7)

	expected, err := maprenderer.RenderIsometricTile(layeredNodeAccessor, layeredColorResolver, p1, p2, nil)
	assert.NoError(t, err)

	img, err := maprenderer.RenderIsometricTile(layeredNodeAccessor, layeredColorResolver, p2, p1, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, img)
}

func TestRenderMapTransparentColor(t *testing.T) {
	// only the water layer is in range, the color has to be written unmodified
	from := types.NewPos(0, 1, 0)
	to := types.NewPos(3, 3, 3)

	img, err := maprenderer.RenderMap(layeredNodeAccessor, layeredColorResolver, from, to, nil)
	assert.NoError(t, err)
	assert.Equal(t, color.NRGBA{R: 0, G: 0, B: 200, A: 100}, img.NRGBAAt(1, 1))
}

func TestRenderIsometricTransparencyComposites(t *testing.T) {
	opts := &maprenderer.IsoRenderOpts{CubeLen: 8, EnableTransparency: true}
	img, err := maprenderer.RenderIsometric(layeredNodeAccessor, layeredColorResolver, types.NewPos(0, 0, 0), types.NewPos(3, 1, 3), opts)
	assert.NoError(t, err)

	rgba := img.(*image.RGBA)
	b := rgba.Bounds()
	c := rgba.RGBAAt(b.Dx()/2, b.Dy()/2)

	// water on top of stone: the stone has to shine through and the result is opaque
	assert.Equal(t, uint8(255), c.A)
	assert.Greater(t, c.R, uint8(0))
	assert.Less(t, c.B, uint8(200))
	assert.Greater(t, c.B, c.R)
}

func TestDrawIsoCubeComposites(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	opaque := color.NRGBA{R: 200, G: 0, B: 0, A: 255}
	half := color.NRGBA{R: 0, G: 0, B: 200, A: 128}

	assert.NoError(t, maprenderer.DrawIsoCube(img, 8, 0, 0, opaque, opaque, opaque))
	assert.NoError(t, maprenderer.DrawIsoCube(img, 8, 0, 0, half, half, half))

	// reference: image/draw "over" compositing
	expected := image.NewRGBA(image.Rect(0, 0, 1, 1))
	expected.Set(0, 0, opaque)
	draw.Draw(expected, expected.Bounds(), image.NewUniform(half), image.Point{}, draw.Over)

	assert.Equal(t, expected.RGBAAt(0, 0), img.RGBAAt(3, 4))
}

func TestDrawIsoCubeTransparentOnEmpty(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	water := color.NRGBA{R: 0, G: 0, B: 200, A: 100}
	assert.NoError(t, maprenderer.DrawIsoCube(img, 8, 0, 0, water, water, water))

	// stored premultiplied, reads back as the original color
	nrgba := color.NRGBAModel.Convert(img.At(3, 4)).(color.NRGBA)
	assert.Equal(t, uint8(100), nrgba.A)
	assert.InDelta(t, 200, int(nrgba.B), 2)
	assert.Equal(t, uint8(0), nrgba.R)
}

func TestBlendColorAlpha(t *testing.T) {
	// opaque background: result is opaque
	c := maprenderer.BlendColor(&color.RGBA{R: 100, G: 100, B: 100, A: 255}, &color.RGBA{B: 200, A: 100}, 1)
	assert.Equal(t, uint8(255), c.A)

	// translucent over translucent: alpha accumulates (a_f + a_b*(1-a_f)) instead of max(a_f, a_b)
	c = maprenderer.BlendColor(&color.RGBA{R: 200, A: 100}, &color.RGBA{B: 200, A: 100}, 1)
	assert.Equal(t, uint8(161), c.A)
	assert.Greater(t, c.R, uint8(0))
	assert.Greater(t, c.B, c.R)

	// fully transparent background does not tint the foreground
	c = maprenderer.BlendColor(&color.RGBA{}, &color.RGBA{B: 200, A: 100}, 1)
	assert.Equal(t, color.RGBA{B: 200, A: 100}, *c)
}

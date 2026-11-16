package maprenderer_test

import (
	"image/color"
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

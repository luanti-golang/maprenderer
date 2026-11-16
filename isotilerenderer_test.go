package maprenderer_test

import (
	"image/png"
	"os"
	"testing"

	"github.com/luanti-golang/colormapping"
	"github.com/luanti-golang/maprenderer"
	"github.com/luanti-golang/types"
	"github.com/stretchr/testify/assert"
)

func TestIsoTileRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	from := types.NewPos(16, 0, 16)
	to := types.NewPos(31, 62, 31)
	opts := &maprenderer.IsoTileRenderOpts{
		CubeLen: 16,
	}

	img, err := maprenderer.RenderIsometricTile(m.GetNode, cm.GetColor, from, to, opts)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	f, err := os.OpenFile("output/iso-tile-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}

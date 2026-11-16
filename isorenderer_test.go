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

func TestIsoRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	from := types.NewPos(0, 0, 0)
	to := types.NewPos(15, 30, 15)
	opts := &maprenderer.IsoRenderOpts{
		CubeLen: 16,
	}

	img, err := maprenderer.RenderIsometric(m.GetNode, cm.GetColor, from, to, opts)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	f, err := os.OpenFile("output/iso-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}

package maprenderer_test

import (
	"testing"

	"github.com/luanti-golang/colormapping"
	"github.com/luanti-golang/maprenderer"
	"github.com/luanti-golang/types"
)

func setupBenchmark(b *testing.B) (*Map, *colormapping.ColorMapping) {
	m := NewMap()
	if err := m.Load("testdata/map.csv"); err != nil {
		b.Fatal(err)
	}
	cm := colormapping.NewColorMapping()
	if err := cm.LoadDefaults(); err != nil {
		b.Fatal(err)
	}
	return m, cm
}

var benchFrom, benchTo = types.NewPos(0, -20, 0), types.NewPos(100, 50, 100)

func BenchmarkRenderMap(b *testing.B) {
	m, cm := setupBenchmark(b)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := maprenderer.RenderMap(m.GetNode, cm.GetColor, benchFrom, benchTo, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderIsometric(b *testing.B) {
	m, cm := setupBenchmark(b)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := maprenderer.RenderIsometric(m.GetNode, cm.GetColor, benchFrom, benchTo, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderIsometricTile(b *testing.B) {
	m, cm := setupBenchmark(b)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := maprenderer.RenderIsometricTile(m.GetNode, cm.GetColor, types.NewPos(0, -20, 0), types.NewPos(47, 50, 47), nil); err != nil {
			b.Fatal(err)
		}
	}
}

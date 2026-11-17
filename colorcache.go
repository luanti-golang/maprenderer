package maprenderer

import (
	"image/color"

	"github.com/luanti-golang/types"
)

type colorKey struct {
	name   string
	param2 int
}

// cachedColorResolver memoizes the results of the given resolver, including misses (nil).
// Resolvers may be expensive for unknown names (fallback matching), and the same few names are queried
// over and over during a render. The returned resolver is not safe for concurrent use.
func cachedColorResolver(cr types.ColorResolver) types.ColorResolver {
	cache := map[colorKey]*color.RGBA{}
	return func(name string, param2 int) *color.RGBA {
		key := colorKey{name: name, param2: param2}
		c, found := cache[key]
		if !found {
			c = cr(name, param2)
			cache[key] = c
		}
		return c
	}
}

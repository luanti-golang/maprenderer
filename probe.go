package maprenderer

import (
	"fmt"

	"github.com/luanti-golang/types"
)

func Probe(min, max, pos, ipos *types.Pos, na types.NodeAccessor, cr types.ColorResolver, skip_alpha bool) ([]*NodeWithColor, error) {
	return probeInto([]*NodeWithColor{}, min, max, pos, ipos, na, cr, skip_alpha)
}

// probeInto works like Probe but appends the found nodes to the given slice
func probeInto(nodes []*NodeWithColor, min, max, pos, ipos *types.Pos, na types.NodeAccessor, cr types.ColorResolver, skip_alpha bool) ([]*NodeWithColor, error) {
	cpos := pos
	for cpos.IsWithin(min, max) {
		node, err := na(cpos)
		if err != nil {
			return nil, fmt.Errorf("getNode error @ %s: %v", cpos, err)
		}

		if isRenderable(node) {
			c := cr(node.Name, node.Param2)
			if c != nil {
				nodes = append(nodes, &NodeWithColor{
					Node:  node,
					Color: c,
				})

				if c.A == 255 || skip_alpha {
					// solid color or skip-param set
					break
				}
			}
		}
		cpos = cpos.Add(ipos)
	}

	return nodes, nil
}

// isRenderable returns false for nodes that are never drawn
func isRenderable(node *types.Node) bool {
	return node != nil && node.Name != "air" && node.Name != "ignore"
}

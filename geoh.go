package geoh

import (
	"fmt"
	"strings"

	"github.com/hsrodriguez/go-geoh/internal/tools"
	"github.com/mmcloughlin/geohash"
)

// This function will convert the given geojson to a list of geohashes
// with the specified precision.
// It will return a list of geohashes as strings.
func Geohashes(geojson string, precision uint, start_precision uint) []string {
	geohashes := []string{}

	mp := tools.GetMultiPolygon(geojson)

	p := start_precision
	if p > precision {
		p = precision
	}

	// Get the center geohash of the multipolygon
	center := tools.GetCenterGeohash(mp, p)
	// Add the center geohash to the list
	geohashes = append(geohashes, center)

	// Get the neighbors of the center geohash
	neighbors := geohash.Neighbors(center)
	// Add the neighbors to the list
	for _, neighbor := range neighbors {
		// Check if the neighbor intersects with the multipolygon
		if tools.Intersect(tools.GetGeohashPolygon(neighbor), mp) {
			// If it does, add it to the list
			geohashes = append(geohashes, neighbor)
		}
	}

	for p < precision {
		innerGeohashes := []string{}
		for _, gh := range geohashes {
			innerGeohashes = append(innerGeohashes, tools.GetInnerGeohashes(gh)...)
		}
		geohashes = geohashes[:0] // Clear the list to avoid duplicates
		for _, gh := range innerGeohashes {
			// Check if the inner geohash intersects with the multipolygon
			if tools.Intersect(tools.GetGeohashPolygon(gh), mp) {
				// If it does, add it to the list
				geohashes = append(geohashes, gh)
			}
		}
		p++
	}

	fmt.Println(strings.Join(geohashes, ","))

	return geohashes
}

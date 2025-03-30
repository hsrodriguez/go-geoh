package tools

import (
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
)

const geohashBase = "0123456789bcdefghjkmnpqrstuvwxyz"

// This function will calculate the center geohash of a given multipolygon
// with the specified precision.
// It will return the geohash as a string.
func GetCenterGeohash(mp orb.MultiPolygon, precision uint) string {
	centroid := MultiPolygonCentroid(mp)
	gh := geohash.EncodeWithPrecision(centroid[1], centroid[0], precision) // Note: geohash takes latitude, longitude
	return gh
}

func GetInnerGeohashes(geohash string) []string {
	result := []string{}
	for _, r := range geohashBase {
		gh := geohash + string(r)
		result = append(result, gh)
	}
	return result
}

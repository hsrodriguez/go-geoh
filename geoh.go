package geoh

import (
	"fmt"

	"github.com/hsrodriguez/go-geoh/internal/tools"
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	gj "github.com/paulmach/orb/geojson" // Import the geojson subpackage
)

// This function will convert the given geojson to a list of geohashes
// with the specified precision.
// It will return a list of geohashes as strings.
func Geohashes(geojson string, precision uint, start_precision uint) []string {
	var geohashes []string

	mp := getMultiPolygon(geojson)

	p := start_precision
	if p > precision {
		p = precision
	}

	// Get the center geohash of the multipolygon
	center := getCenterGeohash(mp, p)
	// Add the center geohash to the list
	geohashes = append(geohashes, center)

	// Get the neighbors of the center geohash
	neighbors := geohash.Neighbors(center)
	// Add the neighbors to the list
	geohashes = append(geohashes, neighbors...)

	fmt.Println(center)

	return geohashes
}

// This function will calculate the center geohash of a given multipolygon
// with the specified precision.
// It will return the geohash as a string.
func getCenterGeohash(mp orb.MultiPolygon, precision uint) string {
	centroid := tools.MultiPolygonCentroid(mp)
	gh := geohash.EncodeWithPrecision(centroid[1], centroid[0], precision) // Note: geohash takes latitude, longitude
	return gh
}

// This function will extract the multipolygon from the given geojson
// It will return the multipolygon as an orb.MultiPolygon.
func getMultiPolygon(geojson string) orb.MultiPolygon {
	var multiPolygon orb.MultiPolygon

	featureCollection, _ := gj.UnmarshalFeatureCollection([]byte(geojson))
	if featureCollection != nil {
		for _, feature := range featureCollection.Features {
			if feature.Geometry != nil {
				switch geom := feature.Geometry.(type) {
				case orb.Polygon:
					multiPolygon = append(multiPolygon, geom)
				case orb.MultiPolygon:
					multiPolygon = append(multiPolygon, geom...)
				}
			}
		}
	}
	feature, _ := gj.UnmarshalFeature([]byte(geojson))
	if feature != nil {
		if feature.Geometry != nil {
			switch geom := feature.Geometry.(type) {
			case orb.Polygon:
				multiPolygon = append(multiPolygon, geom)
			case orb.MultiPolygon:
				multiPolygon = append(multiPolygon, geom...)
			}
		}
	}

	return multiPolygon
}

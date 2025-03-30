package tools

import (
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"

	gj "github.com/paulmach/orb/geojson" // Import the geojson subpackage
)

// This function will extract the multipolygon from the given geojson
// It will return the multipolygon as an orb.MultiPolygon.
func GetMultiPolygon(geojson string) orb.MultiPolygon {
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

// This function will get the polygon of the given geohash
// It will return the polygon as an orb.Polygon.
func GetGeohashPolygon(gh string) orb.Polygon {
	geohashPolygon := geohash.BoundingBox(gh)
	result := orb.Polygon{}
	result = append(result, orb.Ring{
		{geohashPolygon.MinLng, geohashPolygon.MinLat},
		{geohashPolygon.MaxLng, geohashPolygon.MinLat},
		{geohashPolygon.MaxLng, geohashPolygon.MaxLat},
		{geohashPolygon.MinLng, geohashPolygon.MaxLat},
		{geohashPolygon.MinLng, geohashPolygon.MinLat},
	})
	return result
}

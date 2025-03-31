package geoh

import (
	gh "github.com/mmcloughlin/geohash"
	gj "github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

// This function will convert the given geojson to a list of geohashes
// with the specified precision.
// It will return a list of geohashes as strings.
func Geohashes(geojson string, precision uint, start_precision uint) []string {
	geohashes := []string{}

	//Parse the geojson to a feature collection
	if fc, err := gj.Parse(geojson, nil); err != nil {
		//Panics if the geojsoon is invalid
		panic(err)
	} else {

		//define the start precision
		p := start_precision
		if p > precision {
			p = precision
		}

		//find the center geohash of the start precision
		centerPoint := fc.Center()
		center := gh.EncodeWithPrecision(centerPoint.Y, centerPoint.X, p)
		// add the center geohash to the results
		geohashes = append(geohashes, center)
		// find the neighbors of the center geohash
		neighbors := gh.Neighbors(center)
		for _, neighbor := range neighbors {
			// get the neighbor geohash polygon
			p := getGeohashPolygon(neighbor)
			// find if the neighbor geohash intersects with the feature collection
			if fc.Intersects(p) {
				//if it does, append to the current geohashes
				geohashes = append(geohashes, neighbor)
			}
		}

		for p < precision {
			inner := []string{}
			//find the inner geohashes of each geohash
			for _, gh := range geohashes {
				inner = append(inner, getInnerGeohashes(gh)...)
			}
			//clean the geohashes to avoid duplicates
			geohashes = geohashes[:0]
			//iterate through the inner geohashes
			for _, gh := range inner {
				//convert the geohash to a polygon
				p := getGeohashPolygon(gh)
				//find if the geohash intersects with the feature collection
				if fc.Intersects(p) {
					//if it does, append to the current geohashes
					geohashes = append(geohashes, gh)
				}
			}
			//increase the precision
			p++
		}
	}

	return geohashes
}

// This function converts the geohash to a polygon
// Returns a geojson.Polygon
func getGeohashPolygon(geohash string) *gj.Polygon {
	//get the bounding box
	bbox := gh.BoundingBox(geohash)
	//create the polygon points
	points := []geometry.Point{
		{X: bbox.MinLng, Y: bbox.MinLat},
		{X: bbox.MaxLng, Y: bbox.MinLat},
		{X: bbox.MaxLng, Y: bbox.MaxLat},
		{X: bbox.MinLng, Y: bbox.MaxLat},
		{X: bbox.MinLng, Y: bbox.MinLat},
	}

	//generate the poly
	poly := geometry.NewPoly(points, nil, nil)
	//return the polygon
	return gj.NewPolygon(poly)
}

// This function gets the geohashes inside a geohash
// Returns an array of geohashes in string
func getInnerGeohashes(geohash string) []string {
	const geohashBase = "0123456789bcdefghjkmnpqrstuvwxyz"
	result := []string{}
	//iterate through the base
	for _, r := range geohashBase {
		gh := geohash + string(r)
		result = append(result, gh)
	}
	return result
}

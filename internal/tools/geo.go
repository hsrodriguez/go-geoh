package tools

import (
	"math"

	"github.com/paulmach/orb"
)

func MultiPolygonCentroid(multiPolygon orb.MultiPolygon) orb.Point {
	var lonSum, latSum, areaSum float64

	// Iterate through each polygon in the multipolygon
	for _, polygon := range multiPolygon {
		centroid := PolygonCentroid(polygon) // Calculate centroid for each polygon
		area := polygonArea(polygon)         //compute area of the polygon

		lonSum += centroid[0] * area
		latSum += centroid[1] * area
		areaSum += area
	}

	if areaSum == 0 {
		return multiPolygon[0][0][0] // return the first point if area is zero
	}

	//compute final weighted centroid
	centroidLon := lonSum / areaSum
	centroidLat := latSum / areaSum

	return orb.Point{centroidLon, centroidLat} // correct order: longitude, latitude
}

// This function will calculate the area of a given polygon
// using the spherical triangle area method.
func polygonArea(polygon orb.Polygon) float64 {
	var areaSum float64
	ring := polygon[0]
	n := len(ring)

	for i := 0; i < n-1; i++ {
		areaSum += sphericalTriangleArea(
			ring[0][0], ring[i][1],
			ring[i][0], ring[i][1],
			ring[i+1][0], ring[i+1][1])
	}
	return areaSum
}

// This function will calculate the centroid of a given polygon
// using the spherical triangle area method.
// It will return the centroid as an orb.Point.
func PolygonCentroid(polygon orb.Polygon) orb.Point {
	var lonSum, latSum, areaSum float64

	//get the exterior ring
	ring := polygon[0]
	n := len(ring)

	//Iterate throufh each triangle formed by the first vertex and each pair of consecutive vertices
	for i := 0; i < n-1; i++ {
		//get the three points forming a spherical triangle
		p1 := ring[0] // fist vertex as reference
		p2 := ring[i]
		p3 := ring[i+1]

		// extract the latitude and longitude
		lon1, lat1 := p1[0], p1[1]
		lon2, lat2 := p2[0], p2[1]
		lon3, lat3 := p3[0], p3[1]

		//compute the spherical area of the triangle
		area := sphericalTriangleArea(lon1, lat1, lon2, lat2, lon3, lat3)
		areaSum += area

		//compute area weighted sum of latitude and longitude
		lonSum += area * (lon1 + lon2 + lon3) / 3
		latSum += area * (lat1 + lat2 + lat3) / 3
	}

	if areaSum == 0 {
		return ring[0] // return the first point if area is zero
	}

	//compute final weighted centroid
	centroidLon := lonSum / areaSum
	centroidLat := latSum / areaSum

	return orb.Point{centroidLon, centroidLat} // correct order: longitude, latitude
}

// This function will calculate the spherical triangle area
// using the haversine formula and l'huillier's formula.
func sphericalTriangleArea(lon1, lat1, lon2, lat2, lon3, lat3 float64) float64 {
	//Convert degrees to radians
	lon1, lat1 = toRadians(lon1), toRadians(lat1)
	lon2, lat2 = toRadians(lon2), toRadians(lat2)
	lon3, lat3 = toRadians(lon3), toRadians(lat3)

	//Calculate the spherical excess using the haversine formula
	a := haversineDistance(lon2, lat2, lon3, lat3)
	b := haversineDistance(lon1, lat1, lon3, lat3)
	c := haversineDistance(lon1, lat1, lon2, lat2)

	//calculate semi-perimeter
	s := (a + b + c) / 2

	//apply l'huillier's formula to calculate the area
	tanE := math.Tan(s/2) * math.Tan((s-a)/2) * math.Tan((s-b)/2) * math.Tan((s-c)/2)
	if tanE <= 0 {
		return 0
	}
	area := 4 * math.Atan(math.Sqrt(tanE))
	return area
}

// This function will convert degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// This function will calculate the haversine distance
// between two points on the sphere.
func haversineDistance(lon1, lat1, lon2, lat2 float64) float64 {
	dLon := (lon2 - lon1) / 2
	dLat := (lat2 - lat1) / 2
	a := math.Sin(dLat)*math.Sin(dLat) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon)*math.Sin(dLon)
	return 2 * math.Asin(math.Sqrt(a))
}

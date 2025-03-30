package tools

import "github.com/paulmach/orb"

// Intersect checks if a polygon intersects with any polygon in a multipolygon
// using the spherical triangle area method.
// It returns true if there is an intersection, otherwise false.
func Intersect(polygon orb.Polygon, multipolygon orb.MultiPolygon) bool {
	// Iterate through each polygon in the multipolygon
	for _, p := range multipolygon {
		if polygonIntersects(polygon, p) {
			return true
		}
	}
	return false
}

func polygonIntersects(polygon1, polygon2 orb.Polygon) bool {
	// Iterate through each ring in both polygons
	for _, ring1 := range polygon1 {
		for _, ring2 := range polygon2 {
			// Check if any point in ring1 is inside ring2
			for _, point1 := range ring1 {
				if pointInRing(ring2, point1) {
					return true
				}
			}

			// Check if any point in ring2 is inside ring1
			for _, point2 := range ring2 {
				if pointInRing(ring1, point2) {
					return true
				}
			}

			// Check if any edge in ring1 intersects with any edge in ring2
			for i := 0; i < len(ring1)-1; i++ {
				for j := 0; j < len(ring2)-1; j++ {
					if edgesIntersect(ring1[i], ring1[i+1], ring2[j], ring2[j+1]) {
						return true
					}
				}
			}
		}
	}
	return false
}

// pointInRing checks if a point is inside a ring using the ray-casting algorithm.
// It returns true if the point is inside the ring, otherwise false.
func pointInRing(ring orb.Ring, point orb.Point) bool {
	result := false
	for i := 0; i < len(ring)-1; i++ {
		edgeStart := ring[i]
		edgeEnd := ring[i+1]
		if (edgeStart[1] > point[1]) != (edgeEnd[1] > point[1]) &&
			(point[0] < (edgeEnd[0]-edgeStart[0])*(point[1]-edgeStart[1])/(edgeEnd[1]-edgeStart[1])+edgeStart[0]) {
			result = !result
		}
	}
	return result
}

// edgesIntersect checks if two edges intersect using the determinant method.
// It returns true if they intersect, otherwise false.
func edgesIntersect(p1, p2, p3, p4 orb.Point) bool {
	det := (p1[1]-p2[1])*(p3[0]-p4[0]) - (p1[0]-p2[0])*(p3[1]-p4[1])
	if det == 0 {
		return false // lines are parallel
	}
	t1 := ((p1[0]-p3[0])*(p3[1]-p4[1]) - (p1[1]-p3[1])*(p3[0]-p4[0])) / det
	t2 := ((p1[0]-p3[0])*(p1[1]-p2[1]) - (p1[1]-p3[1])*(p1[0]-p2[0])) / det

	return t1 >= 0 && t1 <= 1 && t2 >= 0 && t2 <= 1
}
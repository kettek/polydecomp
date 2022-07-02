package polydecomp

import (
	"math"
)

type Polygon[T float32 | float64] [][2]T

func pointIntersection[T float32 | float64](p1, p2, q1, q2 [2]T) [2]T {
	var p [2]T
	a1 := p2[1] - p1[1]
	b1 := p1[0] - p2[0]
	c1 := a1*p1[0] + b1*p1[1]
	a2 := q2[1] - q1[1]
	b2 := q1[0] - q2[0]
	c2 := a2*q1[0] + b2*q1[1]
	det := a1*b2 - a2*b1
	// Lines are not parallel
	if !(math.Abs(float64(det)) <= 1e-8) {
		p[0] = (b2*c1 - b1*c2) / det
		p[1] = (a1*c2 - a2*c1) / det
	}
	return p
}

func pointArea[T float32 | float64](a, b, c [2]T) T {
	return ((b[0] - a[0]) * (c[1] - a[1])) - ((c[0] - a[0]) * (b[1] - a[1]))
}
func pointLeft[T float32 | float64](a, b, c [2]T) bool {
	return pointArea(a, b, c) > 0
}
func pointLeftOn[T float32 | float64](a, b, c [2]T) bool {
	return pointArea(a, b, c) >= 0
}
func pointRight[T float32 | float64](a, b, c [2]T) bool {
	return pointArea(a, b, c) < 0
}
func pointRightOn[T float32 | float64](a, b, c [2]T) bool {
	return pointArea(a, b, c) <= 0
}
func pointSquareDistance[T float32 | float64](a, b [2]T) T {
	dx := b[0] - a[0]
	dy := b[1] - a[1]
	return dx*dx + dy*dy
}

// CCW ensures a polygon is in counter-clockwise ordering.
func (p Polygon[T]) CCW() {
	br := 0

	// Find bottom right point.
	for i := 1; i < len(p); i++ {
		if p[i][1] < p[br][1] || (p[i][1] == p[br][1] && p[i][0] > p[br][0]) {
			br = i
		}
	}

	// Reverse if clockwise.
	if !pointLeft(p.at(br-1), p.at(br), p.at(br+1)) {
		for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
			p[i], p[j] = p[j], p[i]
		}
	}
}

func (p Polygon[T]) isReflex(i int) bool {
	return pointRight(p.at(i-1), p.at(i), p.at(i+1))
}

func (p Polygon[T]) at(i int) [2]T {
	if i < 0 {
		return p[i%len(p)+len(p)]
	}
	return p[i%len(p)]
}

func (polygon Polygon[T]) decompose(polys *[]Polygon[T]) {
	var upperInt, lowerInt, p [2]T
	var upperDist, lowerDist, d, closestDist T
	var upperIndex, lowerIndex, closestIndex int
	var lowerPoly, upperPoly Polygon[T]

	insert := func(dst *Polygon[T], src Polygon[T], start, end int) {
		for i := start; i < end; i++ {
			*dst = append(*dst, src[i])
		}
	}

	for i := 0; i < len(polygon); i++ {
		if polygon.isReflex(i) {
			// Uh... this is probably a bad idea, but math.Max[Type] does not exist....
			upperDist, lowerDist = 0, 0
			upperDistChecked, lowerDistChecked := false, false

			for j := 0; j < len(polygon); j++ {
				if pointLeft(polygon.at(i-1), polygon.at(i), polygon.at(j)) && pointRightOn(polygon.at(i-1), polygon.at(i), polygon.at(j-1)) { // If the line intersects with an edge.
					p = pointIntersection(polygon.at(i-1), polygon.at(i), polygon.at(j), polygon.at(j-1)) // Get the intersection point.
					if pointRight(polygon.at(i+1), polygon.at(i), p) {                                    // Ensure it is inside the polygon.
						d = pointSquareDistance(polygon[i], p)
						if d < lowerDist || !lowerDistChecked { // Keep only the closest intersection.
							lowerDistChecked = true
							lowerDist = d
							lowerInt = p
							lowerIndex = j
						}
					}
				}
				if pointLeft(polygon.at(i+1), polygon.at(i), polygon.at(j+1)) && pointRightOn(polygon.at(i+1), polygon.at(i), polygon.at(j)) {
					p = pointIntersection(polygon.at(i+1), polygon.at(i), polygon.at(j), polygon.at(j+1))
					if pointLeft(polygon.at(i-1), polygon.at(i), p) {
						d = pointSquareDistance(polygon[i], p)
						if d < upperDist || !upperDistChecked {
							upperDistChecked = true
							upperDist = d
							upperInt = p
							upperIndex = j
						}
					}
				}
			}

			// If there are no vertices to connect to, choose a point in the middle.
			if lowerIndex == (upperIndex+1)%len(polygon) {
				p[0] = (lowerInt[0] + upperInt[0]) / 2
				p[1] = (lowerInt[1] + upperInt[1]) / 2

				if i < upperIndex {
					insert(&lowerPoly, polygon, i, upperIndex+1)
					lowerPoly = append(lowerPoly, p)
					upperPoly = append(upperPoly, p)
					if lowerIndex != 0 {
						insert(&upperPoly, polygon, lowerIndex, len(polygon))
					}
					insert(&upperPoly, polygon, 0, i+1)
				} else {
					if i != 0 {
						insert(&lowerPoly, polygon, i, len(polygon))
					}
					insert(&lowerPoly, polygon, 0, upperIndex+1)
					lowerPoly = append(lowerPoly, p)
					upperPoly = append(upperPoly, p)
					insert(&upperPoly, polygon, lowerIndex, i+1)
				}
			} else {
				if lowerIndex > upperIndex {
					upperIndex += len(polygon)
				}

				closestDist = 0
				closestDistChecked := false

				for j := lowerIndex; j <= upperIndex; j++ {
					if pointLeftOn(polygon.at(i-1), polygon.at(i), polygon.at(j)) &&
						pointRightOn(polygon.at(i+1), polygon.at(i), polygon.at(j)) {
						d := pointSquareDistance(polygon.at(i), polygon.at(j))
						if d < closestDist || !closestDistChecked {
							closestDistChecked = true
							closestDist = d
							closestIndex = j % len(polygon)
						}
					}
				}

				if i < closestIndex {
					insert(&lowerPoly, polygon, i, closestIndex+1)
					if closestIndex != 0 {
						insert(&upperPoly, polygon, closestIndex, len(polygon))
					}
					insert(&upperPoly, polygon, 0, i+1)
				} else {
					if i != 0 {
						insert(&lowerPoly, polygon, i, len(polygon))
					}
					insert(&lowerPoly, polygon, 0, closestIndex+1)
					insert(&upperPoly, polygon, closestIndex, i+1)
				}
			}

			// Solve smallest poly first.
			if len(lowerPoly) < len(upperPoly) {
				lowerPoly.decompose(polys)
				upperPoly.decompose(polys)
			} else {
				upperPoly.decompose(polys)
				lowerPoly.decompose(polys)
			}
			return
		}
	}

	*polys = append(*polys, polygon)
	return
}

// Decompose decomposes a polygonal shape into convex polygons if needed.
func (p Polygon[T]) Decompose() []Polygon[T] {
	var polys []Polygon[T]
	p.decompose(&polys)
	return polys
}

# polydecomp
This library implements Mark Bayazit's polygon decomposition algorithm, turning a concave polygon into one or more convex polygons.

It requires Go 1.18 due to the use of generics to allow both float32 and float64 types.

Polygons must not be self-intersecting.

## Why
I needed it for generating 2D terrain geometry from arbitrarly shaped polygons and could not find one written in go.

## Basic usage

```
package main

import "github.com/kettek/polydecomp"

func main() {
	poly := polydecomp.Polygon[float64]{
		{-100, 100},
		{-100, 0},
		{100, 0},
		{100, 100},
		{50, 50},
	}

	// Ensure it is counter-clockwise.
	poly.CCW()

	// Decompose it.
	polys := poly.Decompose(math.MaxFloat64)

	// Print 'em out.
	fmt.Println(poly)  // [[-100 100] [-100 0] [100 0] [100 100] [50 50]]
	fmt.Println(polys) // [[[100 0] [100 100] [50 50]] [[50 50] [-100 100] [-100 0] [100 0]]]
}

```


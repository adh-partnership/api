/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package geo

import (
	"runtime"
	"sync"
)

type Point struct {
	Y float64 // Lat
	X float64 // Lon
}

type Polygon struct {
	Points []Point
}

type BoundingBox struct {
	BottomLeft Point
	TopRight   Point
}

func PointInPolygon(pt Point, poly Polygon) bool {
	bbox := GetBoundingBox(poly)
	if !PointInBoundingBox(pt, bbox) {
		return false
	}

	nverts := len(poly.Points)
	intersect := false

	verts := poly.Points
	j := 0

	for i := 1; i < nverts; i++ {
		if ((verts[i].Y > pt.Y) != (verts[j].Y > pt.Y)) &&
			(pt.X < (verts[j].X-verts[i].X)*(pt.Y-verts[i].Y)/(verts[j].Y-verts[i].Y)+verts[i].X) {
			intersect = !intersect
		}

		j = i
	}

	return intersect
}

func PointInBoundingBox(pt Point, bb BoundingBox) bool {
	// Check if point is in bounding box

	// Bottom Left is the smallest and x and y value
	// Top Right is the largest x and y value
	return pt.X < bb.TopRight.X && pt.X > bb.BottomLeft.X &&
		pt.Y < bb.TopRight.Y && pt.Y > bb.BottomLeft.Y
}

func GetBoundingBox(poly Polygon) BoundingBox {
	var maxX, maxY, minX, minY float64

	for i := 0; i < len(poly.Points); i++ {
		side := poly.Points[i]

		if side.X > maxX || maxX == 0.0 {
			maxX = side.X
		}
		if side.Y > maxY || maxY == 0.0 {
			maxY = side.Y
		}
		if side.X < minX || minX == 0.0 {
			minX = side.X
		}
		if side.Y < minY || minY == 0.0 {
			minY = side.Y
		}
	}

	return BoundingBox{
		BottomLeft: Point{X: minX, Y: minY},
		TopRight:   Point{X: maxX, Y: maxY},
	}
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func PointInPolygonParallel(pts []Point, poly Polygon, numcores int) []Point {
	MAXPROCS := MaxParallelism()
	runtime.GOMAXPROCS(MAXPROCS)

	if numcores > MAXPROCS {
		numcores = MAXPROCS
	}

	start := 0
	inside := []Point{}

	var m sync.Mutex
	var wg sync.WaitGroup
	wg.Add(numcores)

	for i := 1; i <= numcores; i++ {
		size := (len(pts) / numcores) * i
		batch := pts[start:size]

		go func(batch []Point) {
			defer wg.Done()
			for j := 0; j < len(batch); j++ {
				pt := batch[j]
				if PointInPolygon(pt, poly) {
					m.Lock()
					inside = append(inside, pt)
					m.Unlock()
				}
			}
		}(batch)
		start = size
	}

	wg.Wait()
	return inside
}

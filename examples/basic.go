package main

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/AeroNotix/quadtree"
)

func main() {
	n := 10000
	canvasSize := 500
	q := quadtree.QuadTree{
		MaxPointsPerNode: 12,
		BoundingBox:      image.Rect(0, 0, canvasSize, canvasSize),
	}

	for x := 0; x < n; x++ {
		point := image.Point{rand.Intn(canvasSize), rand.Intn(canvasSize)}
		q.InsertPoint(point)
	}
	fmt.Println(q.Draw("qtree.png"))
}

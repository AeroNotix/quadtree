package quadtree

import (
	"bufio"
	"image"

	"code.google.com/p/draw2d/draw2d"
	"image/png"

	"log"
	"os"
)

type QuadTree struct {
	MaxPointsPerNode int
	points           []image.Point
	BoundingBox      image.Rectangle
	BLeft, TLeft     *QuadTree
	BRight, TRight   *QuadTree
	hasChildren      bool
}

func (q *QuadTree) InsertPoint(p image.Point) {
	if q.hasChildren {
		if p.In(q.BLeft.BoundingBox) {
			q.BLeft.InsertPoint(p)
			return
		}
		if p.In(q.BRight.BoundingBox) {
			q.BRight.InsertPoint(p)
			return
		}
		if p.In(q.TLeft.BoundingBox) {
			q.TLeft.InsertPoint(p)
			return
		}
		if p.In(q.TRight.BoundingBox) {
			q.TRight.InsertPoint(p)
			return
		}
	}
	q.points = append(q.points, p)
	if !q.acceptingPoints() && !q.hasChildren {
		q.rebalance()
	}
}

func (q *QuadTree) acceptingPoints() bool {
	return len(q.points) < q.MaxPointsPerNode
}

func (q *QuadTree) LowerLeft() image.Rectangle {
	return image.Rect(
		q.BoundingBox.Min.X,
		q.BoundingBox.Min.Y,
		(q.BoundingBox.Min.X+q.BoundingBox.Max.X)/2,
		(q.BoundingBox.Min.Y+q.BoundingBox.Max.Y)/2,
	)
}

func (q *QuadTree) UpperLeft() image.Rectangle {
	p := image.Point{
		q.BoundingBox.Max.X,
		q.BoundingBox.Min.Y,
	}
	return image.Rect(
		(p.X+q.BoundingBox.Min.X)/2,
		(p.Y+q.BoundingBox.Min.Y)/2,
		(p.X+q.BoundingBox.Max.X)/2,
		(p.Y+q.BoundingBox.Max.Y)/2,
	)
}

func (q *QuadTree) UpperRight() image.Rectangle {
	return image.Rect(
		(q.BoundingBox.Min.X+q.BoundingBox.Max.X)/2,
		(q.BoundingBox.Min.Y+q.BoundingBox.Max.Y)/2,
		q.BoundingBox.Max.X,
		q.BoundingBox.Max.Y,
	)
}

func (q *QuadTree) LowerRight() image.Rectangle {
	Min := q.UpperLeft().Min
	return image.Rect(
		Min.X,
		Min.Y,
		q.BoundingBox.Max.X,
		q.BoundingBox.Max.Y/2,
	)
}

func (q *QuadTree) rebalance() {
	if q.acceptingPoints() {
		return
	}
	q.TLeft = &QuadTree{
		MaxPointsPerNode: q.MaxPointsPerNode,
		BoundingBox:      q.UpperLeft(),
	}
	q.BLeft = &QuadTree{
		MaxPointsPerNode: q.MaxPointsPerNode,
		BoundingBox:      q.LowerLeft(),
	}
	q.TRight = &QuadTree{
		MaxPointsPerNode: q.MaxPointsPerNode,
		BoundingBox:      q.UpperRight(),
	}
	q.BRight = &QuadTree{
		MaxPointsPerNode: q.MaxPointsPerNode,
		BoundingBox:      q.LowerRight(),
	}
	for _, p := range q.points {
		if p.In(q.BLeft.BoundingBox) {
			q.BLeft.InsertPoint(p)
			continue
		}
		if p.In(q.BRight.BoundingBox) {
			q.BRight.InsertPoint(p)
			continue
		}
		if p.In(q.TLeft.BoundingBox) {
			q.TLeft.InsertPoint(p)
			continue
		}
		if p.In(q.TRight.BoundingBox) {
			q.TRight.InsertPoint(p)
			continue

		}
	}
	q.points = nil
	q.hasChildren = true
}

func (q QuadTree) Walk() []QuadTree {
	if q.acceptingPoints() {
		return []QuadTree{q}
	}
	if q.hasChildren {
		var nodes []QuadTree
		nodes = append(nodes, q.BLeft.Walk()...)
		nodes = append(nodes, q.TLeft.Walk()...)
		nodes = append(nodes, q.TRight.Walk()...)
		return append(nodes, q.BRight.Walk()...)
	}
	return []QuadTree{}
}

func saveToPngFile(filePath string, m image.Image) error {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return b.Flush()
}

func (q *QuadTree) drawOnContext(gc *draw2d.ImageGraphicContext) {
	gc.MoveTo(float64(q.BoundingBox.Max.X), float64(q.BoundingBox.Max.Y))
	gc.LineTo(float64(q.BoundingBox.Max.X), float64(q.BoundingBox.Max.X))
	gc.LineTo(float64(q.BoundingBox.Min.X), float64(q.BoundingBox.Min.Y))
	gc.LineTo(float64(q.BoundingBox.Max.X), float64(q.BoundingBox.Min.Y))
	gc.Stroke()
}

func (q *QuadTree) Draw(fpath string) error {
	Nodes := q.Walk()
	img := image.NewRGBA(image.Rect(0, 0, q.BoundingBox.Max.X, q.BoundingBox.Max.Y))
	gc := draw2d.NewGraphicContext(img)

	for _, node := range Nodes {
		node.drawOnContext(gc)
	}

	return saveToPngFile(fpath, img)
}

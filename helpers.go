package MarsRover

import "math"

// Distance returns the distance between the point and the given one.
func (p Point) Distance(other Point) float64 {
	return math.Pow(float64(other.X)-float64(p.X), 2) +
		math.Pow(float64(other.Y)-float64(p.Y), 2)
}

// Equal returns whether the point is equal to the given point.
func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

package tokyoclient

import (
	"math"
)

// CalculateAngle calculates the angle to the target
func CalculateAngle(p1, p2 Point) float64 {
	vectorAB := []float64{p2.X - p1.X, p2.Y - p1.Y}

	angle := math.Atan2(vectorAB[1], vectorAB[0])

	// Ensure the angle is in the range [0, 2π)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}

// DistanceBetween calculates the distance to the target
func DistanceBetween(p1, p2 Point) float64 {
	vectorAB := []float64{p2.X - p1.X, p2.Y - p1.Y}

	return math.Sqrt(vectorAB[0]*vectorAB[0] + vectorAB[1]*vectorAB[1])
}

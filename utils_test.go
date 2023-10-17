package tokyoclient

import (
	"math"
	"testing"
)

func TestCalculateAngle(t *testing.T) {
	testCases := map[string]struct {
		p1, p2   Point
		expected float64
	}{
		"Horizontal Line": {
			p1:       Point{X: 0, Y: 0},
			p2:       Point{X: 1, Y: 0},
			expected: 0.0,
		},
		"Vertical Line": {
			p1:       Point{X: 0, Y: 0},
			p2:       Point{X: 0, Y: 1},
			expected: math.Pi / 2.0,
		},
		"Diagonal Line": {
			p1:       Point{X: 1, Y: 1},
			p2:       Point{X: 0, Y: 0},
			expected: -3*math.Pi/4.0 + 2*math.Pi,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := CalculateAngle(tc.p1, tc.p2)
			if math.Abs(result-tc.expected) > 1e-9 {
				t.Errorf("Expected angle %f, but got %f", tc.expected, result)
			}
		})
	}
}

func TestCalculateDistance(t *testing.T) {
	testCases := map[string]struct {
		p1, p2   Point
		expected float64
	}{
		"In Horizontal Line": {
			p1:       Point{X: 0, Y: 0},
			p2:       Point{X: 1, Y: 0},
			expected: 1.0,
		},
		"In Vertical Line": {
			p1:       Point{X: 0, Y: 0},
			p2:       Point{X: 0, Y: 1},
			expected: 1.0,
		},
		"In Diagonal Line": {
			p1:       Point{X: 1, Y: 1},
			p2:       Point{X: 4, Y: 5},
			expected: 5.0,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := DistanceBetween(tc.p1, tc.p2)
			if result != tc.expected {
				t.Errorf("Expected distance %f, but got %f", tc.expected, result)
			}
		})
	}
}

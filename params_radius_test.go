package ingv

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRadiusDistanceInKm(t *testing.T) {
	records := [][]float64{
		//        lat1        lon1        lat2        lon2        distance/km  error margin
		// distance Naples-Rome
		[]float64{40.8539041, 14.1642017, 41.9099533, 12.3711904, 190.4, 1},
		// distance Naples-Pompeii
		[]float64{40.8539041, 14.1642017, 40.7466178, 14.4730826, 28.66, 1},
		// distance Naples-London
		[]float64{40.8539041, 14.1642017, 51.5269879, -0.7253499, 1645.7, 2},
	}
	for _, record := range records {
		computedDistance := DistanceInKm(record[0], record[1], record[2], record[3])
		expectedDistance := record[4]
		errorMargin := record[5]
		assert.LessOrEqual(t, math.Abs(computedDistance-expectedDistance), errorMargin)
	}
}

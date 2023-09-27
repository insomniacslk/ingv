package ingv

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRadiusDistanceInKm(t *testing.T) {
	type record struct {
		origin, dest         Loc
		distanceKm, marginKm float64
	}
	records := []record{
		// distance Naples-Rome
		record{
			origin:     Loc{Lat: 40.8539041, Lon: 14.1642017},
			dest:       Loc{Lat: 41.9099533, Lon: 12.3711904},
			distanceKm: 190.4,
			marginKm:   1,
		},
		// distance Naples-Pompeii
		record{
			origin:     Loc{Lat: 40.8539041, Lon: 14.1642017},
			dest:       Loc{Lat: 40.7466178, Lon: 14.4730826},
			distanceKm: 28.66,
			marginKm:   1},
		// distance Naples-London
		record{
			origin:     Loc{Lat: 40.8539041, Lon: 14.1642017},
			dest:       Loc{Lat: 51.5269879, Lon: -0.7253499},
			distanceKm: 1645.7,
			marginKm:   2,
		},
	}
	for _, r := range records {
		computedDistance := DistanceInKm(r.origin, r.dest)
		assert.LessOrEqual(t, math.Abs(computedDistance-r.distanceKm), r.marginKm)
	}
}

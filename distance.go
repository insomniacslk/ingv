package ingv

import "math"

type Loc struct {
	Lat float64
	Lon float64
}

// DistanceInKm computes the distance in KM between two sets of coordinates
// using the haversine formula, that is, assuming a spherical Earth (the error
// margin is ~0.3%).
// This function works around what appears to be a broken min/max radius
// computation in the API.
//
// Warning: using this function is different than using With{Min,Max}RadiusKm,
// because the filtering happens after the search, so you have to handle that
// manually.
func DistanceInKm(origin, dest Loc) float64 {
	latRad := deg2rad(origin.Lat - dest.Lat)
	lonRad := deg2rad(origin.Lon - dest.Lon)
	a := (math.Sin(latRad/2)*math.Sin(latRad/2) +
		math.Cos(deg2rad(origin.Lat))*
			math.Cos(deg2rad(dest.Lat))*math.Sin(lonRad/2)*
			math.Sin(lonRad/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	// 6371 is the Earth's mean radius in km
	return 6371 * c
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

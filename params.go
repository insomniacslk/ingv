package ingv

import (
	"math"
	"time"
)

type Param func(r *Request)

func WithEventID(id int) Param {
	return func(r *Request) {
		r.EventID = &id
	}
}

func WithOriginID(id int) Param {
	return func(r *Request) {
		r.OriginID = &id
	}
}

func WithMagnitudeID(id int) Param {
	return func(r *Request) {
		r.MagnitudeID = &id
	}
}

func WithFocalMechanismID(id int) Param {
	return func(r *Request) {
		r.FocalMechanismID = &id
	}
}

func WithMinLat(l float64) Param {
	return func(r *Request) {
		r.MinLat = &l
	}
}

func WithMaxLat(l float64) Param {
	return func(r *Request) {
		r.MaxLat = &l
	}
}

func WithMinLon(l float64) Param {
	return func(r *Request) {
		r.MinLon = &l
	}
}

func WithMaxLon(l float64) Param {
	return func(r *Request) {
		r.MaxLon = &l
	}
}

func WithLat(l float64) Param {
	return func(r *Request) {
		r.Lon = &l
	}
}

func WithLon(l float64) Param {
	return func(r *Request) {
		r.Lon = &l
	}
}

// WARNING: the following four functions (min/max radius in degrees/km)
// appear to be broken in the current API implementation.
// See DistanceInKm for a workaround.

func WithMaxRadius(rad float64) Param {
	return func(r *Request) {
		r.MaxRadius = &rad
	}
}

func WithMaxRadiusKm(rad float64) Param {
	return func(r *Request) {
		r.MaxRadiusKm = &rad
	}
}

func WithMinRadius(rad float64) Param {
	return func(r *Request) {
		r.MinRadius = &rad
	}
}

func WithMinRadiusKm(rad float64) Param {
	return func(r *Request) {
		r.MinRadiusKm = &rad
	}
}

func WithMinDepth(d float64) Param {
	return func(r *Request) {
		r.MinDepth = &d
	}
}

func WithMaxDepth(d float64) Param {
	return func(r *Request) {
		r.MaxDepth = &d
	}
}

func WithStartTime(t time.Time) Param {
	return func(r *Request) {
		r.StartTime = &t
	}
}

func WithEndTime(t time.Time) Param {
	return func(r *Request) {
		r.EndTime = &t
	}
}

func WithMinMag(m float64) Param {
	return func(r *Request) {
		r.MinMag = &m
	}
}

func WithMaxMag(m float64) Param {
	return func(r *Request) {
		r.MaxMag = &m
	}
}

func WithMagnitudeType(m int) Param {
	return func(r *Request) {
		r.MagnitudeType = &m
	}
}

func WithLimit(l int) Param {
	return func(r *Request) {
		r.Limit = &l
	}
}

func WithOffset(o int) Param {
	return func(r *Request) {
		r.Offset = &o
	}
}

func WithOrderBy(o string) Param {
	return func(r *Request) {
		r.OrderBy = &o
	}
}

func WithIncludeAllMagnitude(i bool) Param {
	return func(r *Request) {
		r.IncludeAllMagnitude = &i
	}
}

func WithIncludeArrivals(i bool) Param {
	return func(r *Request) {
		r.IncludeArrivals = &i
	}
}

func WithIncludeAllOrigins(i bool) Param {
	return func(r *Request) {
		r.IncludeAllOrigins = &i
	}
}

func WithIncludeAllStationMagnitudes(i bool) Param {
	return func(r *Request) {
		r.IncludeAllStationMagnitudes = &i
	}
}

func WithFormat(f string) Param {
	return func(r *Request) {
		r.Format = &f
	}
}

// DistanceInKm computes the distance in KM between two sets of coordinates
// using the haversine formula, that is, assuming a spherical Earth (the error
// margin is ~0.3%).
// This function works around what appears to be a broken min/max radius
// computation in the API.
//
// Warning: using this function is different than using With{Min,Max}RadiusKm,
// because the filtering happens after the search, so it may yield fewer results.
func DistanceInKm(lat1, lon1, lat2, lon2 float64) float64 {
	latRad := deg2rad(lat1 - lat2)
	lonRad := deg2rad(lon1 - lon2)
	a := (math.Sin(latRad/2)*math.Sin(latRad/2) +
		math.Cos(deg2rad(lat1))*
			math.Cos(deg2rad(lat2))*math.Sin(lonRad/2)*
			math.Sin(lonRad/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	// 6371 is the Earth's mean radius in km
	return 6371 * c
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

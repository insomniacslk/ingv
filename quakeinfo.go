package ingv

import "time"

type QuakeInfo struct {
	EventID           int
	Time              time.Time
	Latitude          float64
	Longitude         float64
	DepthInKm         float64
	Author            string
	Catalog           string
	Contributor       string
	ContributorID     int
	MagType           string
	Magnitude         float64
	MagAuthor         string
	EventLocationName string
	EventType         string
}

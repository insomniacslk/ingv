package main

import (
	"fmt"
	"log"
	"time"

	"github.com/insomniacslk/ingv"
	"github.com/spf13/pflag"
)

var (
	flagMinMagnitude = pflag.Float64P("min-magnitude", "m", 2.0, "Minimum magnitude to report")
	flagMaxMagnitude = pflag.Float64P("max-magnitude", "M", 10.0, "Maximum magnitude to report")
	flagStartTime    = pflag.DurationP("start-time", "s", -24*time.Hour, "Start time, relative (FIXME: add support for absolute time)")
	flagEndTime      = pflag.DurationP("end-time", "e", 0, "End time, relative (FIXME: add support for absolute time)")
	flagLatitude     = pflag.Float64P("latitude", "l", 0.0, "Latitude. If unspecified do not restrict by location. If specified, longitude is mandatory")
	flagLongitude    = pflag.Float64P("longitude", "L", 0.0, "Longitude. If unspecified do not restrict by location. If specified, latitude is mandatory")
	flagMaxRadius    = pflag.Float64P("max-radius", "R", 0.0, "Max radius in degrees from lat/lon. If unspecified, no radius limit is set. If specified, lat/lon are mandatory")
	flagMinRadius    = pflag.Float64P("min-radius", "r", 0.0, "Min radius in degrees from lat/lon. If unspecified, no radius limit is set. If specified, lat/lon are mandatory")
	flagLimit        = pflag.IntP("limit", "i", 10, "Max number of results to show")
	flagOrderBy      = pflag.StringP("order-by", "b", "time", "Order by time or magnitude")
)

// it appears that the radius is not computed properly by the API, so this flag
// controls the use of a workaround. Note that the results with the workaround
// might be fewer than expected.
const buggyRadiusAPI = true

func main() {
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	var params []ingv.Param
	params = append(params, ingv.WithMinMag(*flagMinMagnitude))
	params = append(params, ingv.WithMaxMag(*flagMaxMagnitude))
	now := time.Now()
	start := now.Add(*flagStartTime)
	params = append(params, ingv.WithStartTime(start))
	end := now.Add(*flagEndTime)
	params = append(params, ingv.WithEndTime(end))
	latLonSpecified := false
	if pflag.CommandLine.Changed("latitude") && pflag.CommandLine.Changed("longitude") {
		latLonSpecified = true
		params = append(params, ingv.WithLat(*flagLatitude))
		params = append(params, ingv.WithLon(*flagLongitude))
	} else if pflag.CommandLine.Changed("latitude") || pflag.CommandLine.Changed("longitude") {
		log.Fatalf("Both latitude and longitude must be specified")
	}
	if !buggyRadiusAPI {
		if pflag.CommandLine.Changed("max-radius") {
			if !latLonSpecified {
				log.Fatalf("latitude and longitude must be specified when using max-radius")
			}
			params = append(params, ingv.WithMaxRadius(*flagMaxRadius))
		}
		if pflag.CommandLine.Changed("min-radius") {
			if !latLonSpecified {
				log.Fatalf("latitude and longitude must be specified when using min-radius")
			}
			params = append(params, ingv.WithMinRadius(*flagMinRadius))
		}
	}
	params = append(params, ingv.WithLimit(*flagLimit))
	params = append(params, ingv.WithOrderBy(*flagOrderBy))
	// only "text" format is supported for now
	params = append(params, ingv.WithFormat("text"))

	records, err := ingv.Get(params...)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	if len(records) == 0 {
		fmt.Printf("No earthquakes found with the specified parameters\n")
	}
	idx := 0
	for _, rec := range records {
		if buggyRadiusAPI {
			if pflag.CommandLine.Changed("max-radius") {
				// work around an API bug. However this is not equivalent, as it may return
				// fewer results than expected.
				if ingv.DistanceInKm(*flagLatitude, *flagLongitude, rec.Latitude, rec.Longitude) > *flagMaxRadius {
					continue
				}
			}
			if pflag.CommandLine.Changed("min-radius") {
				// work around an API bug. However this is not equivalent, as it may return
				// fewer results than expected.
				if ingv.DistanceInKm(*flagLatitude, *flagLongitude, rec.Latitude, rec.Longitude) < *flagMinRadius {
					continue
				}
			}
		}
		fmt.Printf("%d) %s\n    Location: %s\n    Magnitude: %.1f\n    Map: https://www.google.com/maps/search/%f,%f/@%f,%f\n    Details: https://terremoti.ingv.it/event/%d for details\n", idx+1, rec.Time, rec.EventLocationName, rec.Magnitude, rec.Latitude, rec.Longitude, rec.Latitude, rec.Longitude, rec.EventID)
		idx++
	}

}

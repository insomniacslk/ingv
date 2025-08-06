package ingv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// time format used by the INGV API
const TimeFormat = "2006-01-02T15:04:05"

var TimeLocation = time.UTC

// API: https://webservices.ingv.it/swagger-ui/dist/?url=https://ingv.github.io/openapi/fdsnws/event/0.0.1/event.yaml
// example: http://webservices.ingv.it/fdsnws/event/1/query?starttime=2012-05-29T00:00:00&endtime=2012-05-29T23:59:59&format=text

type Request struct {
	EventID                     *int       `name:"eventid"`
	OriginID                    *int       `name:"originid"`
	MagnitudeID                 *int       `name:"magnitudeid"`
	FocalMechanismID            *int       `name:"focalmechanismid"`
	MinLat                      *float64   `name:"minlat"`
	MaxLat                      *float64   `name:"maxlat"`
	MinLon                      *float64   `name:"minlon"`
	MaxLon                      *float64   `name:"maxlon"`
	Lat                         *float64   `name:"lat"`
	Lon                         *float64   `name:"lon"`
	MaxRadius                   *float64   `name:"maxradius"`
	MaxRadiusKm                 *float64   `name:"maxradiuskm"`
	MinRadius                   *float64   `name:"minradius"`
	MinRadiusKm                 *float64   `name:"minradiuskm"`
	MinDepth                    *float64   `name:"mindepth"`
	MaxDepth                    *float64   `name:"maxdepth"`
	StartTime                   *time.Time `name:"starttime"`
	EndTime                     *time.Time `name:"endtime"`
	MinMag                      *float64   `name:"minmag"`
	MaxMag                      *float64   `name:"maxmag"`
	MagnitudeType               *int       `name:"magnitudetype"`
	Limit                       *int       `name:"limit"`
	Offset                      *int       `name:"offset"`
	OrderBy                     *string    `name:"orderby"`
	IncludeAllMagnitude         *bool      `name:"includeallmagnitude"`
	IncludeArrivals             *bool      `name:"includearrivals"`
	IncludeAllOrigins           *bool      `name:"includeallorigins"`
	IncludeAllStationMagnitudes *bool      `name:"includeallstationmagnitudes"`
	Format                      *string    `name:"format"`
}

func (r *Request) ToValues() url.Values {
	values := url.Values{}
	e := reflect.ValueOf(r).Elem()
	for idx := 0; idx < e.NumField(); idx++ {
		field := e.Field(idx)

		if field.IsNil() {
			// field is not set, do not add it to the params
			continue
		}
		name := e.Type().Field(idx).Tag.Get("name")
		intf := field.Interface()
		var value string

		switch v := intf.(type) {
		case *int:
			value = strconv.FormatInt(int64(*v), 10)
		case *float64:
			value = strconv.FormatFloat(*v, 'f', -1, 64)
		case *bool:
			value = strconv.FormatBool(*v)
		case *string:
			value = string(*v)
		case *time.Time:
			value = (*v).Format(TimeFormat)
		default:
			// if this happens, the type of a field of the Request struct is
			// not handled in this switch, whoops. Just add it
			log.Fatalf("Unhandled type %T for field `%s`: bug? See source code comments", v, name)
		}
		values.Add(name, value)
	}
	return values
}

func Get(params ...Param) ([]QuakeInfo, error) {
	r := Request{}
	for _, p := range params {
		p(&r)
	}

	if r.Format == nil || *r.Format != "text" {
		return nil, fmt.Errorf("only `format=text` is supported in the request")
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "webservices.ingv.it",
		Path:     "/fdsnws/event/1/query",
		RawQuery: r.ToValues().Encode(),
	}
	log.Printf("Connecting to %s", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// per-IP rate limit
	// TODO handle rate limit
	ratelimit := resp.Header.Get("X-RateLimit-Limit")
	ratelimitReset := resp.Header.Get("X-RateLimit-Reset")
	log.Printf("Rate limit: %s, reset: %s", ratelimit, ratelimitReset)
	if resp.StatusCode == 204 {
		// HTTP 204 No Content
		return nil, nil
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("expected HTTP 200 OK, got %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return parseQuakeInfoCSV(body)
}

func parseQuakeInfoCSV(body []byte) ([]QuakeInfo, error) {

	errorAt := func(r *csv.Reader, name string, idx int, err error) error {
		line, column := r.FieldPos(idx)
		return fmt.Errorf("failed to parse field %s (line: %d, column: %d, index: 0): %w", name, line, column, err)
	}

	parseIntField := func(r *csv.Reader, fields []string, name string, fieldIdx int) (int, error) {
		i, err := strconv.ParseInt(fields[fieldIdx], 10, 64)
		if err != nil {
			return 0, errorAt(r, name, fieldIdx, err)
		}
		return int(i), nil
	}

	parseFloatField := func(r *csv.Reader, fields []string, name string, fieldIdx int) (float64, error) {
		f, err := strconv.ParseFloat(fields[fieldIdx], 64)
		if err != nil {
			return 0, errorAt(r, name, fieldIdx, err)
		}
		return f, nil
	}

	parseTimeField := func(r *csv.Reader, fields []string, name string, fieldIdx int) (*time.Time, error) {
		t, err := time.ParseInLocation(TimeFormat, fields[1], TimeLocation)
		if err != nil {
			return nil, errorAt(r, name, fieldIdx, err)
		}
		return &t, nil
	}
	ret := make([]QuakeInfo, 0)
	r := csv.NewReader(bytes.NewReader(body))
	r.Comma = '|'
	r.Comment = '#'
	for fieldIdx := 0; ; fieldIdx++ {
		fields, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var qi QuakeInfo
		// event ID
		i, err := parseIntField(r, fields, "EventID", 0)
		if err != nil {
			return nil, err
		}
		qi.EventID = i
		// time
		t, err := parseTimeField(r, fields, "Time", 1)
		if err != nil {
			return nil, err
		}
		qi.Time = *t
		// latitude
		f, err := parseFloatField(r, fields, "Latitude", 2)
		if err != nil {
			return nil, err
		}
		qi.Latitude = f
		// longitude
		f, err = parseFloatField(r, fields, "Longitude", 3)
		if err != nil {
			return nil, err
		}
		qi.Longitude = f
		// DepthInKm
		f, err = parseFloatField(r, fields, "Depth/Km", 4)
		if err != nil {
			return nil, err
		}
		qi.DepthInKm = f
		// author
		qi.Author = fields[5]
		// catalog
		qi.Catalog = fields[6]
		// contributor
		qi.Contributor = fields[7]
		// contibutor ID
		if fields[8] != "" {
			// non-mandatory field: parse only if non-empty
			i, err = parseIntField(r, fields, "ContributorID", 8)
			if err != nil {
				return nil, err
			}
			qi.ContributorID = i
		}
		// mag type
		qi.MagType = fields[9]
		// magnitude
		f, err = parseFloatField(r, fields, "Magnitude", 10)
		if err != nil {
			return nil, err
		}
		qi.Magnitude = f
		// mag author
		qi.MagAuthor = fields[11]
		// event location name
		qi.EventLocationName = fields[12]
		// event type
		qi.EventType = fields[13]

		// add the record to the slice
		ret = append(ret, qi)
	}
	return ret, nil
}

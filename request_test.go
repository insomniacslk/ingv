package ingv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The FDSNWS API takes naked timestamps with no zone and reads them as UTC, so
// a time.Time carried in a non-UTC location has to be CONVERTED, not merely
// formatted. Getting this wrong is silent: the request succeeds, the API
// answers correctly, and the answer is empty because the window asked for has
// not happened yet.
func TestToValuesConvertsTimeToUTC(t *testing.T) {
	// +02:00, i.e. what an Italian caller has for half the year.
	cest := time.FixedZone("CEST", 2*60*60)
	start := time.Date(2026, 9, 12, 20, 55, 1, 0, cest)
	end := start.Add(time.Hour)

	r := Request{}
	WithStartTime(start)(&r)
	WithEndTime(end)(&r)
	values := r.ToValues()

	// 20:55:01+02:00 is 18:55:01Z, which is what the API must be told.
	assert.Equal(t, "2026-09-12T18:55:01", values.Get("starttime"))
	assert.Equal(t, "2026-09-12T19:55:01", values.Get("endtime"))
}

// The same instant must serialise identically whatever location it is carried
// in. This is the property the bug broke: the value sent depended on the
// caller's TZ environment variable.
func TestToValuesIsIndependentOfCallerLocation(t *testing.T) {
	instant := time.Date(2026, 9, 12, 18, 55, 1, 0, time.UTC)
	want := "2026-09-12T18:55:01"

	for _, offset := range []int{-8 * 60 * 60, -60 * 60, 0, 2 * 60 * 60, 9 * 60 * 60} {
		r := Request{}
		WithStartTime(instant.In(time.FixedZone("test", offset)))(&r)
		assert.Equal(t, want, r.ToValues().Get("starttime"),
			"offset %d seconds changed the serialised value", offset)
	}
}

// What ToValues writes, parseTimeField must read back as the same instant --
// the two sides share TimeFormat and TimeLocation and must agree on both.
func TestToValuesRoundTripsThroughTimeLocation(t *testing.T) {
	cest := time.FixedZone("CEST", 2*60*60)
	start := time.Date(2026, 9, 12, 20, 55, 1, 0, cest)

	r := Request{}
	WithStartTime(start)(&r)

	parsed, err := time.ParseInLocation(TimeFormat, r.ToValues().Get("starttime"), TimeLocation)
	require.NoError(t, err)
	assert.True(t, parsed.Equal(start), "round trip moved the instant: %s != %s", parsed, start)
}

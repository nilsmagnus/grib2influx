package main

import (
	"os"
	"testing"

	"github.com/nilsmagnus/grib/griblib"
)

func Test_tofluxpoints(t *testing.T) {

	testFile, fileErr := os.Open("testdata/gfs.t00z.pgrb2.2p50.f003")

	if fileErr != nil {
		t.Fatalf("Error opening testfile %v", fileErr)
	}

	messages, gribErr := griblib.ReadMessages(testFile)

	filtered := griblib.Filter(messages, griblib.Options{
		Discipline: 0,
		Category:   0,
		GeoFilter:  griblib.GeoFilter{
			MinLong:7000000,
			MaxLong:11000000,
			MinLat: 57000000,
			MaxLat: 71000000},
	})

	if gribErr != nil {
		t.Fatalf("Could not parse testfile, %v", gribErr)
	}
	fluxies := toInfluxPoints([]griblib.Message{filtered[0]}, 3)

	if len(fluxies) == 0 || len(fluxies) != int(filtered[0].Section3.DataPointCount) {
		t.Errorf("Expected fluxies length to be the same as message.datapointCount, expected %d, was %d",
			messages[0].Section3.DataPointCount,
			len(fluxies))
	}

	if len(fluxies) == 0 || len(fluxies) != len(filtered[0].Section7.Data) {
		t.Errorf("Expected fluxies length to be the same as message.datapointCount, expected %d, was %d",
			messages[0].Section3.DataPointCount,
			len(fluxies))
	}

}

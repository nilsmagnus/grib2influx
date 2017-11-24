package main

import (
	"os"
	"testing"

	"github.com/nilsmagnus/grib/griblib"
)

func Test_tofluxpoints(t *testing.T) {

	testFile, fileErr := os.Open("testdata/gfs.t00z.pgrb2.2p50.f000")

	if fileErr != nil {
		t.Fatalf("Error opening testfile %v", fileErr)
	}

	messages, gribErr := griblib.ReadMessages(testFile, griblib.Options{})

	if gribErr != nil {
		t.Fatalf("Could not parse testfile, %v", gribErr)
	}
	fluxies := toFlux(messages)

	if len(fluxies) == 0 || len(fluxies) != int(messages[0].Section3.DataPointCount) {
		t.Errorf("Expected fluxies length to be the same as message.datapointCount, was %d",
			messages[0].Section3.DataPointCount)
	}

}

func Test_tocoord(t *testing.T) {
	const di = 2500000
	section3 := griblib.Section3{Definition: &griblib.Grid0{
		Di:  di,
		Dj:  di,
		Lo1: 0,
		Lo2: 357500000,
	}}
	coords := toCoords(45, section3)

	if coords.Lon != (45 * di) {
		t.Errorf("Expected lon %d, got %d", 45*di, coords.Lon)
	}

	count2 := 357500000*3/di + 45

	coords2 := toCoords(count2, section3)

	if coords2.Lon != (45 * di) {
		t.Errorf("Expected lon2 %d, got %d", count2*di, coords.Lon)
	}
	if coords2.Lat != (3 * di) {
		t.Errorf("Expected lat2 %d, got %d", 3*di, coords.Lat)
	}

}

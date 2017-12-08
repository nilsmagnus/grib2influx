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

	messages, gribErr := griblib.ReadMessages(testFile, griblib.Options{MaximumNumberOfMessages: 1})

	if gribErr != nil {
		t.Fatalf("Could not parse testfile, %v", gribErr)
	}
	fluxies := toInfluxPoints(messages, 3)

	if len(fluxies) == 0 || len(fluxies) != int(messages[0].Section3.DataPointCount) {
		t.Errorf("Expected fluxies length to be the same as message.datapointCount, expected %d, was %d",
			messages[0].Section3.DataPointCount,
			len(fluxies))
	}

}

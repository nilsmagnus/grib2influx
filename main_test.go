package main

import (
	"testing"

	"github.com/nilsmagnus/grib/griblib"
)

func Test_tocoord(t *testing.T) {
	const di = 2500000
	const la1 = 90000000
	section3 := griblib.Section3{Definition: &griblib.Grid0{
		Di:  di,
		Dj:  di,
		Lo1: 0,
		Lo2: 357500000,
		La1: la1,
		Nj:  144,
		Ni:  73,
	}}
	coordse := toCoords(45, section3)

	if coordse.Lon != (45 * di) {
		t.Errorf("Expected lon %d, got %d", 45*di, coordse.Lon)
	}

	count2 := 144*3 + 45

	coords2 := toCoords(count2, section3)

	if coords2.Lon != (45 * di) {
		t.Errorf("Expected lon2 %d, got %d, index %v", 45*di, coords2.Lon, coords2.Lon/di)
	}

	if coords2.Lat != la1+(3*di) {
		t.Errorf("Expected lat2 %d, got %d", la1+3*di, coords2.Lat)
	}

}

func Test_offset_from_filename(t *testing.T) {
	offset, err := forecastHourFromFileName("aftenpoften101")
	if err != nil {
		t.Errorf("Should be valid format with three trailing digits, %v", err)
	}

	if offset != 101 {
		t.Errorf("Offset should have been 101, was %d", offset)
	}
}

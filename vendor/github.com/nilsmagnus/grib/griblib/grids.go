package griblib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ScaledValue struct {
	Scale uint8  `json:"scale"`
	Value uint32 `json:"value"`
}

type BasicAngle struct {
	BasicAngle    uint32 `json:"basicAngle"`
	BasicAngleSub uint32 `json:"basicAngleSub"`
}

type Grid interface {
	Export() map[string]string
}

func ReadGrid(f io.Reader, templateNumber uint16) (Grid, error) {
	switch templateNumber {
	case 0:
		var grid Grid0
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	case 10:
		var grid Grid10
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	case 20:
		var grid Grid20
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	case 30:
		var grid Grid30
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	case 40:
		var grid Grid40
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	case 90:
		var grid Grid90
		return &grid, binary.Read(f, binary.BigEndian, &grid)
	default:
		var grid Grid90
		return &grid, errors.New(fmt.Sprint("Unknown grid template number ", templateNumber))
	}
}

type GridHeader struct {
	EarthShape      uint8       `json:"earthShape"`
	SphericalRadius ScaledValue `json:"sphericalRadius"`
	MajorAxis       ScaledValue `json:"majorAxis"`
	MinorAxis       ScaledValue `json:"minorAxis"`
}

func (h *GridHeader) Export() (d map[string]string) {
	return map[string]string{
		"earth": ReadEarthShape(int(h.EarthShape)),
	}
}

// Grid Definition Template 3.0: Latitude/longitude (or equidistant cylindrical, or Plate Carree)
type Grid0 struct {
	//Name :=  "Latitude/longitude (or equidistant cylindrical, or Plate Carree) "
	GridHeader
	Ni                          uint32     `json:"ni"`
	Nj                          uint32     `json:"nj"`
	BasicAngle                  BasicAngle `json:"basicAngle"`
	La1                         int32      `json:"la1"`
	Lo1                         int32      `json:"lo1"`
	ResolutionAndComponentFlags uint8      `json:"resolutionAndComponentFlags"`
	La2                         int32      `json:"la2"`
	Lo2                         int32      `json:"lo2"`
	Di                          int32      `json:"di"`
	Dj                          int32      `json:"dj"`
	ScanningMode                uint8      `json:"scanningMode"`
}

func (h *Grid0) Export() map[string]string {
	return map[string]string{
		"earth":         ReadEarthShape(int(h.EarthShape)),
		"ni":            fmt.Sprint(h.Ni),
		"nj":            fmt.Sprint(h.Nj),
		"basicAngle":    fmt.Sprint(h.BasicAngle.BasicAngle),
		"basicAngleSub": fmt.Sprint(h.BasicAngle.BasicAngleSub),
		"la1":           fmt.Sprint(h.La1),
		"lo1":           fmt.Sprint(h.Lo1),
		"la2":           fmt.Sprint(h.La2),
		"lo2":           fmt.Sprint(h.Lo2),
		"di":            fmt.Sprint(h.Di),
		"dj":            fmt.Sprint(h.Dj),
		"scanningMode":  fmt.Sprint(h.ScanningMode),
	}
}

// Grid Definition Template 3.10: Mercator
type Grid10 struct {
	//name :=  "Mercator"
	GridHeader
	Ni                          uint32 `json:"ni"`
	Nj                          int32  `json:"nj"`
	La1                         int32  `json:"la1"`
	Lo1                         int32  `json:"lo1"`
	ResolutionAndComponentFlags uint8  `json:"resolutionAndComponentFlags"`
	Lad                         int32  `json:"lad"`
	La2                         int32  `json:"la2"`
	Lo2                         int32  `json:"lo2"`
	ScanningMode                uint8  `json:"scanningMode"`
	GridOrientation             uint32 `json:"gridOrientation"`
	Di                          int32  `json:"di"`
	Dj                          int32  `json:"dj"`
}

// Grid Definition Template 3.20: Polar stereographic projection
type Grid20 struct {
	//name =  "Polar stereographic projection ";
	GridHeader
	Nx                          uint32 `json:"Nx"`
	Ny                          uint32 `json:"ny"`
	La1                         int32  `json:"na1"`
	Lo1                         int32  `json:"lo1"`
	ResolutionAndComponentFlags uint8  `json:"resolutionAndComponentFlags"`
	Lad                         int32  `json:"lad"`
	Lov                         int32  `json:"lov"`
	Dx                          int32  `json:"dx"`
	Dy                          int32  `json:"dy"`
	ProjectionCenter            uint8  `json:"projectionCenter"`
	ScanningMode                uint8  `json:"scanningMode"`
}

// Grid Definition Template 3.30: Lambert conformal
type Grid30 struct {
	//name =  "Polar stereographic projection ";
	GridHeader
	Nx                          uint32 `json:"nx"`
	Ny                          uint32 `json:"ny"`
	La1                         int32  `json:"la1"`
	Lo1                         int32  `json:"lo1"`
	ResolutionAndComponentFlags uint8  `json:"resolutionAndComponentFlags"`
	Lad                         int32  `json:"lad"`
	Lov                         int32  `json:"lov"`
	Dx                          int32  `json:"dx"`
	Dy                          int32  `json:"dy"`
	ProjectionCenter            uint8  `json:"projectionCenter"`
	ScanningMode                uint8  `json:"scanningMode"`
	Latin1                      uint32 `json:"latin1"`
	Latin2                      uint32 `json:"latin2"`
	LaSouthPole                 uint32 `json:"laSouthPole"`
	LoSouthPole                 uint32 `json:"loSouthPole"`
}

// Grid Definition Template 3.40: Gaussian latitude/longitude
type Grid40 struct {
	//name =  "Gaussian latitude/longitude ";
	GridHeader
	Ni                          uint32 `json:"ni"`
	Nj                          uint32 `json:"nj"`
	BasicAngle                  uint32 `json:"basicAngle"`
	La1                         int32  `json:"la1"`
	Lo1                         int32  `json:"lo1"`
	ResolutionAndComponentFlags uint8  `json:"resolutionAndComponentFlags"`
	La2                         int32  `json:"la2"`
	Lo2                         int32  `json:"lo2"`
	Di                          int32  `json:"di"`
	N                           uint32 `json:"n"`
	ScanningMode                uint8  `json:"scanningMode"`
}

// Grid Definition Template 3.90: Space view perspective or orthographic
// FIXME: implement properly
type Grid90 struct {
	//name =  "Space view perspective or orthographic ";
	GridHeader
	Nx uint32 `json:"nx"`
	Ny uint32 `json:"ny"`
	//BasicAngle                  BasicAngle
	Lap                         int32  `json:"lap"`
	Lop                         int32  `json:"lop"`
	ResolutionAndComponentFlags uint8  `json:"resolutionAndComponentFlags"`
	Dx                          uint32 `json:"dx"`
	Dy                          uint32 `json:"dy"`
	Xp                          uint32 `json:"xp"`
	Yp                          uint32 `json:"yp"`
	ScanningMode                uint8  `json:"scanningMode"`
	Orientation                 uint32 `json:"orientation"`
	Nr                          uint32 `json:"nr"`
	Xo                          uint32 `json:"xo"`
	Yo                          uint32 `json:"yo"`
}

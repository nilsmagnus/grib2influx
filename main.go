package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"reflect"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/nilsmagnus/grib/griblib"
)

type Coords struct {
	Lat int
	Lon int
}

func main() {
	gribFile := flag.String("gribfile", "", "gribfile to import")
	flag.Parse()

	if *gribFile == "" {
		fmt.Print("No gribfile specified!\n\n")
		flag.Usage()
		os.Exit(1)
	}

	file, ferr := os.Open(*gribFile)

	if ferr != nil {
		panic(ferr)
	}

	messages, err := griblib.ReadMessages(file, griblib.Options{})

	if err != nil {
		fmt.Printf("Error reading gribfile: %v", err)
	}

	fmt.Printf("Read %d messages\n", len(messages))

	fluxies := toFlux(messages)

	storeInDatabase(fluxies)

}
func storeInDatabase(points []*client.Point) {
	panic("not implemented")
}

func toFlux(messages []griblib.Message) []*client.Point {
	points := []*client.Point{}
	for _, message := range messages {
		forecastTime := toGoTime(message.Section1.ReferenceTime)

		for counter, data := range message.Section7.Data {
			coords := toCoords(counter, message.Section3)

			dataTypeName := griblib.ReadProductDisciplineParameters(message.Section0.Discipline,
				message.Section4.ProductDefinitionTemplate.ParameterCategory)
			points = append(points, singleDataPoint(data, dataTypeName, forecastTime, coords))
		}
	}
	return points
}
func toCoords(counter int, section3 griblib.Section3) Coords {
	if grid, ok := section3.Definition.(*griblib.Grid0); ok {
		lonCount := int((grid.Lo2 - grid.Lo1) / grid.Di)

		lat := (counter / lonCount) * int(grid.Di)
		lon := (counter % lonCount) * int(grid.Dj)
		return Coords{
			Lat: lat,
			Lon: lon,
		}

	}

	panic(fmt.Sprintf("Unsupported grid-definition %s", reflect.TypeOf(section3.Definition)))
}

func toGoTime(gTime griblib.Time) time.Time {
	return time.Date(int(gTime.Year), time.Month(gTime.Month), int(gTime.Day), int(gTime.Hour), int(gTime.Minute), int(gTime.Second), 0, time.Now().Location())
}

func singleDataPoint(data int64, dataname string, dataTime time.Time, coords Coords) *client.Point {

	fields := map[string]interface{}{
		dataname: data,
	}
	serieNameForCoordinates := fmt.Sprintf("%dx%d", coords.Lat, coords.Lon)
	tags := map[string]string{}
	point, err := client.NewPoint(serieNameForCoordinates, tags, fields, dataTime)
	if err != nil {
		panic(err)
	}
	return point
}

/*
func asInfluxPoints(data ApixuForecast) []*client.Point {

	points := []*client.Point{}
	for _, day := range data.Forecast.Forecastday {
		for _, hour := range day.Hour {
			fields := map[string]interface{}{
				"celsius":        hour.TempC,
				"wind_kph":       hour.WindKph,
				"wind_degree":    hour.WindDegree,
				"percip_mm":      hour.PrecipMm,
				"cloud":          hour.Cloud,
				"humidity":       hour.Humidity,
				"condition_code": hour.Condition.Code,
			}
			point, _ := client.NewPoint(SerieKey(data.Location), map[string]string{}, fields, time.Unix(int64(hour.TimeEpoch), 0))
			points = append(points, point)
		}
	}
	return points
}
*/

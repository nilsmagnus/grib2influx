package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"reflect"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/nilsmagnus/grib/griblib"
	"strconv"
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

	forecastOffsetHour, forecastOffsetErr := forecastHourFromFileName(*gribFile)

	if forecastOffsetErr != nil {
		panic(fmt.Sprintf("Could not parse forecast offset from filename, expected three trailing digits in filename, got %s", *gribFile))
	}

	messages, err := griblib.ReadMessages(file, griblib.Options{})

	if err != nil {
		fmt.Printf("Error reading gribfile: %v", err)
	}

	fmt.Printf("Read %d messages\n", len(messages))

	fluxies := toInfluxPoints(messages, forecastOffsetHour)

	storeInDatabase(fluxies)

}
func forecastHourFromFileName(fileName string) (int, error) {
	v := fileName[len(fileName)-3:]
	return strconv.Atoi(v)
}

func storeInDatabase(points []*client.Point) {
	panic("not implemented")
}

func toInfluxPoints(messages []griblib.Message, forecastOffsetHour int) []*client.Point {
	points := []*client.Point{}
	for _, message := range messages {

		forecastStartTime := toGoTime(message.Section1.ReferenceTime)

		for counter, data := range message.Section7.Data {
			coords := toCoords(counter, message.Section3)

			dataTypeName := griblib.ReadProductDisciplineParameters(message.Section0.Discipline,
				message.Section4.ProductDefinitionTemplate.ParameterCategory)
			points = append(points, singleInfluxDataPoint(data, dataTypeName, forecastStartTime, coords, forecastOffsetHour))
		}
	}
	return points
}
func toCoords(counter int, section3 griblib.Section3) Coords {
	if grid, ok := section3.Definition.(*griblib.Grid0); ok {
		lonCount := int((grid.Lo2 - grid.Lo1) / grid.Di)

		lat := int(grid.La1) - (counter/lonCount)*int(grid.Di)
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

func singleInfluxDataPoint(data int64, dataname string, forecastTime time.Time, coords Coords, offsetHour int) *client.Point {

	fields := map[string]interface{}{
		fmt.Sprintf("%dx%s", offsetHour, dataname): data,
	}
	// serieName is YYYY-mm-dd-hh.latxlon.dataname
	serieName := fmt.Sprintf("%d-%d-%d-%d.%dx%d.%s",
		forecastTime.Year(), forecastTime.Month(), forecastTime.Day(), forecastTime.Hour(),
		coords.Lat, coords.Lon,
		dataname)

	tags := map[string]string{
		"lat":          fmt.Sprintf("%d", coords.Lat),
		"lon":          fmt.Sprintf("%d", coords.Lon),
		"forecastdate": fmt.Sprintf("%d-%d-%d-%2d", forecastTime.Year(), forecastTime.Month(), forecastTime.Day(), forecastTime.Hour()),
		"offsetHour":   fmt.Sprintf("%d", offsetHour),
	}
	point, err := client.NewPoint(serieName, tags, fields, forecastTime)
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

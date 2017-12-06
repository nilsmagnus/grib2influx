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
	"sync"
)

type Coords struct {
	Lat int
	Lon int
}

const (
	defaultInfluxURL  = "http://localhost"
	defaultInfluxPort = 8086
	defaultDatabase   = "forecasts" // apixu data is actuals
)

func main() {
	gribFile, influxConfig := cliArguments()

	if gribFile == "" {
		fmt.Print("No gribfile specified!\n\n")
		flag.Usage()
		os.Exit(1)
	}

	file, ferr := os.Open(gribFile)

	if ferr != nil {
		panic(ferr)
	}

	forecastOffsetHour, forecastOffsetErr := forecastHourFromFileName(gribFile)

	if forecastOffsetErr != nil {
		panic(fmt.Sprintf("Could not parse forecast offset from filename, expected three trailing digits in filename, got %s", gribFile))
	}

	messages, err := griblib.ReadMessages(file, griblib.Options{})

	if err != nil {
		fmt.Printf("Error reading gribfile: %v", err)
	}

	fmt.Printf("Read %d messages\n", len(messages))

	wg  := sync.WaitGroup{}
	wg.Add(len(messages))

	for _, message := range messages {
		go func() {
			influxPoints := toInfluxPoints([]griblib.Message{message}, forecastOffsetHour)
			saveErr := save(influxPoints, influxConfig)
			if saveErr != nil {
				fmt.Printf("Error saving points in message: %v\n", saveErr)
			}
			wg.Done()
			fmt.Print(".")
		}()
	}
	wg.Wait()
	fmt.Print("\n")

}

func cliArguments() (string, ConnectionConfig) {
	gribFile := flag.String("gribfile", "", "gribfile to import")

	influxHost := flag.String("influxHost", defaultInfluxURL, "Hostname for influxdb")
	influxUser := flag.String("influxUser", "", "User for influxdb")
	influxDatabase := flag.String("database", defaultDatabase, "Database to use")
	influxPassword := flag.String("influxPassword", "", "Password for influx")
	influxPort := flag.Int("influxport", defaultInfluxPort, "Port for influxdb")

	flag.Parse()

	config := ConnectionConfig{
		User:     *influxUser,
		Password: *influxPassword,
		Port:     *influxPort,
		Hostname: *influxHost,
		Database: *influxDatabase,
	}
	flag.Parse()

	return *gribFile, config
}

func forecastHourFromFileName(fileName string) (int, error) {
	v := fileName[len(fileName)-3:]
	return strconv.Atoi(v)
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

func singleInfluxDataPoint(data int64, dataname string, forecastTime time.Time, coords Coords, offsetHours int) *client.Point {

	fields := map[string]interface{}{
		fmt.Sprintf("%s", dataname): data,
	}
	// serieName is YYYY-mm-dd-hh.latxlon.dataname
	serieName := fmt.Sprintf("%d-%d-%d-%02d.%dx%d",
		forecastTime.Year(), forecastTime.Month(), forecastTime.Day(), forecastTime.Hour(),
		coords.Lat/10000, coords.Lon/10000)

	tags := map[string]string{
		"lat":          fmt.Sprintf("%d", coords.Lat),
		"lon":          fmt.Sprintf("%d", coords.Lon),
		"forecastdate": fmt.Sprintf("%d-%d-%d-%02d", forecastTime.Year(), forecastTime.Month(), forecastTime.Day(), forecastTime.Hour()),
		"offsetHours":   fmt.Sprintf("%d", offsetHours),
	}
	valueTime := forecastTime.Add(time.Duration(offsetHours) * time.Hour)
	point, err := client.NewPoint(serieName, tags, fields, valueTime)
	if err != nil {
		panic(err)
	}
	return point
}

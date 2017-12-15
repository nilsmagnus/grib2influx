package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"reflect"

	"strconv"

	"github.com/nilsmagnus/grib/griblib"
)

type Coords struct {
	Lat int
	Lon int
}

const (
	defaultInfluxURL  = "http://localhost"
	defaultInfluxPort = 8086
	defaultDatabase   = "forecasts"
)

func main() {
	portNo, gribFile, influxConfig := cliArguments()

	if gribFile == "" {
		fmt.Print("No gribfile specified, starting server mode")
		startServerMode(portNo, influxConfig)
	}

	file, ferr := os.Open(gribFile)

	if ferr != nil {
		panic(ferr)
	}

	forecastOffsetHour, forecastOffsetErr := forecastHourFromFileName(gribFile)

	if forecastOffsetErr != nil {
		panic(fmt.Sprintf("Could not parse forecast offset from filename, expected three trailing digits in filename, got %s", gribFile))
	}

	messages, err := griblib.ReadMessages(file)

	if err != nil {
		fmt.Printf("Error reading gribfile: %v", err)
		panic(err)
	}

	fmt.Printf("Read %d messages\n", len(messages))

	// filter the messages based on what criterias you might have. See griblib github page for examples
	filtered := griblib.Filter(messages, griblib.Options{
		Discipline: 0, //
		Category:   0, // temperature
		// norway+ sweden (ish)
		GeoFilter: griblib.GeoFilter{
			MinLat:  57000000,
			MaxLat:  71000000,
			MinLong: 4400000,
			MaxLong: 32000000,
		},
		Surface: griblib.Surface{
			Type:  100, // isobaric surface
			Value: 100,
		},
	})

	fmt.Printf("%d messages after filtering \n", len(filtered))

	client, err := clientFromConfig(influxConfig)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	for _, message := range filtered {
		influxPoints := toInfluxPoints([]griblib.Message{message}, forecastOffsetHour)
		saveErr := save(influxPoints, client, influxConfig.Database)
		if saveErr != nil {
			fmt.Printf("Error saving points in message: %v\n", saveErr)
		}
		fmt.Print(".")
	}
	fmt.Print("\n")

}

func cliArguments() (portno int, gribfileName string, influxConfig ConnectionConfig) {
	gribFile := flag.String("gribfile", "", "Gribfile to import. If no gribfile specified, start server mode.")

	portNo := flag.Int("port", 8080, "Server port no, if servermode.")
	influxHost := flag.String("influxHost", defaultInfluxURL, "Hostname for influxdb.")
	influxUser := flag.String("influxUser", "", "User for influxdb.")
	influxDatabase := flag.String("database", defaultDatabase, "Database name to use.")
	influxPassword := flag.String("influxPassword", "", "Password for influx.")
	influxPort := flag.Int("influxport", defaultInfluxPort, "Port for influxdb.")

	flag.Parse()

	influxConfig = ConnectionConfig{
		User:     *influxUser,
		Password: *influxPassword,
		Port:     *influxPort,
		Hostname: *influxHost,
		Database: *influxDatabase,
	}
	flag.Parse()

	return *portNo, *gribFile, influxConfig
}

func forecastHourFromFileName(fileName string) (int, error) {
	v := fileName[len(fileName)-3:]
	return strconv.Atoi(v)
}

func toCoords(counter int, section3 griblib.Section3) Coords {
	if grid, ok := section3.Definition.(*griblib.Grid0); ok {
		lonCount := int(grid.Nj)

		lat := int(grid.La1) + (counter/lonCount)*int(grid.Di)
		lon := int(grid.Lo1) + (counter%lonCount)*int(grid.Dj)
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

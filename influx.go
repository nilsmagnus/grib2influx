package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/nilsmagnus/grib/griblib"
	"time"
)

//ConnectionConfig is connection config w/credentials and url
type ConnectionConfig struct {
	Hostname string
	Port     int
	User     string
	Password string
	Database string
}

func clientFromConfig(config ConnectionConfig) (client.Client, error) {
	return client.NewHTTPClient(client.HTTPConfig{
		Addr:               fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		Username:           config.User,
		Password:           config.Password,
		InsecureSkipVerify: true,
	})
}

func save(points []*client.Point, config ConnectionConfig) error {
	c, err := clientFromConfig(config)

	if err != nil {
		return fmt.Errorf("Error creating client %s", err.Error())
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.Database,
		Precision: "s",
	})

	if err != nil {
		return fmt.Errorf("Error creating batchpoints %s", err.Error())
	}

	bp.AddPoints(points)

	return c.Write(bp)
}

func toInfluxPoints(messages []griblib.Message, forecastOffsetHour int) []*client.Point {
	points := []*client.Point{}
	for _, message := range messages {

		forecastStartTime := toGoTime(message.Section1.ReferenceTime)

		dataTypeName := griblib.ReadProductDisciplineParameters(message.Section0.Discipline,
			message.Section4.ProductDefinitionTemplate.ParameterCategory)

		for counter, data := range message.Section7.Data {
			coords := toCoords(counter, message.Section3)

			points = append(points, singleInfluxDataPoint(data, dataTypeName, forecastStartTime, coords, forecastOffsetHour))
		}
	}
	return points
}


func singleInfluxDataPoint(data int64, dataname string, forecastTime time.Time, coords Coords, offsetHours int) *client.Point {

	fields := map[string]interface{}{
		fmt.Sprintf("%s-%06dx%06d", dataname, coords.Lat/10000, coords.Lon/10000): data,
	}

	serieName := fmt.Sprintf("%d-%02d-%02d-%02d",
		forecastTime.Year(), forecastTime.Month(), forecastTime.Day(), forecastTime.Hour())

	valueTime := forecastTime.Add(time.Duration(offsetHours) * time.Hour)
	point, err := client.NewPoint(serieName, map[string]string{}, fields, valueTime)
	if err != nil {
		panic(err)
	}
	return point
}


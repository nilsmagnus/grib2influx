package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
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


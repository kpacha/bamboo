package configuration

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/influxdb/influxdb/client"
)

type Influxdb struct {
	Enabled         bool
	Host            string
	DB              string
	Username        string
	Password        string
	Prefix          string
	sanitizedPrefix string
	Connection      *client.Client
	point           chan client.Point
}

func (i *Influxdb) CreateClient() {
	if !i.Enabled || i.Connection != nil {
		return
	}
	u, err := url.Parse(fmt.Sprintf("http://%s", i.Host))
	if err != nil {
		log.Fatal("Error creating the influxdb url: ", err)
	}

	connectionConf := client.Config{
		URL:      *u,
		Username: i.Username,
		Password: i.Password,
	}

	i.Connection, err = client.NewClient(connectionConf)
	if err != nil {
		log.Fatal("Error connecting to the influxdb store: ", err)
	}
	i.sanitizedPrefix = i.Prefix
	if i.sanitizedPrefix != "" {
		i.sanitizedPrefix = fmt.Sprintf("%s.", i.sanitizedPrefix)
	}
	i.point = make(chan client.Point, 1000)

	go i.client()
}

func (i *Influxdb) Save(measurement string, value interface{}) {
	i.point <- client.Point{
		Measurement: fmt.Sprintf("%s%s", i.sanitizedPrefix, measurement),
		Fields: map[string]interface{}{
			"value": value,
		},
		Time: time.Now(),
	}
}

func (i *Influxdb) client() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		if err := i.flush(); err != nil {
			log.Fatal("Error flushing to the influxdb store: ", err)
		}
	}
}

func (i *Influxdb) flush() error {
	bps := client.BatchPoints{
		Points:          i.getPoints(),
		Database:        i.DB,
		RetentionPolicy: "default",
	}
	_, err := i.Connection.Write(bps)
	return err
}

func (i *Influxdb) getPoints() []client.Point {
	var buffer []client.Point
	for {
		select {
		case p := <-i.point:
			buffer = append(buffer, p)
		default:
			break
		}
	}
	return buffer
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/peterhellberg/ruuvitag"
)

func setup(ctx context.Context) context.Context {
	d, err := dev.DefaultDevice()
	if err != nil {
		panic(err)
	}
	ble.SetDefaultDevice(d)

	return ble.WithSigHandler(context.WithCancel(ctx))
}

func main() {
	ctx := setup(context.Background())

	ble.Scan(ctx, true, handler, filter)
}

func handler(a ble.Advertisement) {

	go func(a ble.Advertisement) {
		db := pg.Connect(&pg.Options{
			Addr:     ":5432",
			User:     "ruuvi",
			Password: "ruuvi",
			Database: "ruuvi",
		})

		defer db.Close()

		raw, err := ruuvitag.ParseRAWv2(a.ManufacturerData())
		if err == nil {
			fmt.Printf("[%s] RSSI: %3d: %+v\n", a.Addr(), a.RSSI(), raw)

			event := &Event{
				Id:             uuid.New(),
				Address:        a.Addr().String(),
				Rssi:           a.RSSI(),
				Temperature:    raw.Temperature,
				Humidity:       raw.Humidity,
				Pressure:       raw.Pressure,
				Acceleration_x: raw.Acceleration.X,
				Acceleration_y: raw.Acceleration.Y,
				Acceleration_z: raw.Acceleration.Z,
				Battery:        raw.Battery,
				Movement:       raw.Movement,
				Sequence:       raw.Sequence,
				Timestamp:      int32(time.Now().Unix()),
			}

			err = db.Insert(event)
			if err != nil {
				panic(err)
			}

		}
	}(a)

}

func filter(a ble.Advertisement) bool {
	return ruuvitag.IsRAWv2(a.ManufacturerData())
}

type Event struct {
	Id             uuid.UUID
	Address        string
	Rssi           int
	Temperature    float64
	Humidity       float64
	Pressure       uint32
	Acceleration_x int16
	Acceleration_y int16
	Acceleration_z int16
	Battery        uint16
	Movement       uint8
	Sequence       uint16
	Timestamp      int32
}

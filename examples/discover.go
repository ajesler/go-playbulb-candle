package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
)

const (
	candleService = "ff02"
)

var (
	sigs        = make(chan os.Signal, 1)
	done        = make(chan bool, 1)
	peripherals = make([]string, 0)
)

func containsString(ss []string, s string) bool {
	for _, e := range ss {
		if s == e {
			return true
		}
	}
	return false
}

func onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		// only look for devices that advertise a Playbulb candle service
		d.Scan([]gatt.UUID{gatt.MustParseUUID(candleService)}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	knownPeripheral := containsString(peripherals, p.ID())

	if !knownPeripheral {
		fmt.Printf("Found '%s' with ID '%s'\n", a.LocalName, p.ID())

		peripherals = append(peripherals, p.ID())
	}
}

func main() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		fmt.Println("error starting:", err)
	}

	d.Handle(
		gatt.PeripheralDiscovered(onPeripheralDiscovered),
	)

	d.Init(onStateChanged)

	fmt.Println("Scanning for Playbulb Candles...")

	<-done
}

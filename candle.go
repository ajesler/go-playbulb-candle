package playbulb

import (
	"errors"
	"fmt"
	"github.com/currantlabs/gatt"
	"github.com/currantlabs/gatt/examples/option"
	"strings"
	"time"
)

const (
	candleServiceUUID = "ff02"
	effectCharUUID    = "fffb"
	colourCharUUID    = "fffc"
)

type Candle struct {
	id               string
	per              gatt.Peripheral
	colourChar       *gatt.Characteristic
	effectChar       *gatt.Characteristic
	doneChannel      chan struct{}
	connectedChannel chan bool
	connected        bool
}

func NewCandle(id string) *Candle {
	p := Candle{
		id:               id,
		doneChannel:      make(chan struct{}),
		connectedChannel: make(chan bool),
		connected:        false,
	}
	return &p
}

func (p *Candle) Off() {
	c := NewColour(0, 0, 0, 0)
	e := p.solidColourEffect(c)

	p.SetEffect(e)
}

func (p *Candle) SetColour(c *Colour) error {
	e := p.solidColourEffect(c)
	payload := p.effectPayload(e)

	err := p.per.WriteCharacteristic(p.effectChar, payload, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *Candle) ReadColour() ([]byte, error) {
	currentColour, err := p.per.ReadCharacteristic(p.colourChar)

	if err != nil {
		fmt.Println("Failed to read characteristic:", err)
		return nil, err
	}

	return currentColour, nil
}

func (p *Candle) SetEffect(e *Effect) error {
	payload := p.effectPayload(e)

	err := p.per.WriteCharacteristic(p.effectChar, payload, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *Candle) ReadEffect() ([]byte, error) {
	currentEffect, err := p.per.ReadCharacteristic(p.effectChar)
	if err != nil {
		fmt.Println("Failed to read characteristic:", err)
		return nil, err
	}
	return currentEffect, nil
}

func (p *Candle) Connect() error {
	if p.connected {
		return nil
	}

	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		return err
	}

	d.Handle(
		gatt.PeripheralDiscovered(p.onPeripheralDiscovered),
		gatt.PeripheralConnected(p.onPeripheralConnected),
		gatt.PeripheralDisconnected(p.onPeripheralDisconnected),
	)

	d.Init(p.onStateChanged)

	select {
	case c := <-p.connectedChannel:
		p.connected = c
		return nil
	case <-time.After(time.Second * 5):
		p.connected = false
		return errors.New(fmt.Sprintf("Failed to connect to candle %s before the five second timeout expired", p.id))
	}
}

func (p *Candle) Disconnect() {
	if p.connected {
		p.per.Device().CancelConnection(p.per)
	}
}

func (p *Candle) colourPayload(c *Colour) []byte {
	return []byte{c.Brightness(), c.R(), c.G(), c.B()}
}

func (p *Candle) effectPayload(e *Effect) []byte {
	return []byte{
		e.Colour().Brightness(), e.Colour().R(), e.Colour().G(), e.Colour().B(),
		e.Mode(), 0, e.Speed(), 0}
}

func (p *Candle) solidColourEffect(c *Colour) *Effect {
	return NewEffect(SOLID, c, 0)
}

func (p *Candle) onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func (p *Candle) onPeripheralDiscovered(per gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if strings.ToUpper(per.ID()) != strings.ToUpper(p.id) {
		return
	}

	per.Device().StopScanning()
	per.Device().Connect(per)
}

func (p *Candle) onPeripheralConnected(per gatt.Peripheral, err error) {
	services, err := per.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range services {
		if s.UUID().String() != candleServiceUUID {
			continue
		}

		cs, err := per.DiscoverCharacteristics(nil, s)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			continue
		}

		for _, c := range cs {
			cUUID := c.UUID().String()
			switch {
			case cUUID == effectCharUUID:
				p.effectChar = c
			case cUUID == colourCharUUID:
				p.colourChar = c
			default:
				continue
			}
		}
	}

	p.per = per

	p.connectedChannel <- true
}

func (p *Candle) onPeripheralDisconnected(per gatt.Peripheral, err error) {
	p.connected = false

	p.per = nil
	p.colourChar = nil
	p.effectChar = nil
}

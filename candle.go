package playbulb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
)

const (
	candleServiceUUID = "ff02"
	effectCharUUID    = "fffb"
	colourCharUUID    = "fffc"
	toggleCharUUID    = "2a37"
)

type Candle interface {
	Connect() error
	Disconnect()
	IsConnected() bool

	SetEffect(*Effect) error
	SetColour(*Colour) error
	OnToggle(func(bool)) error
	Off() error
}

type candle struct {
	id               string
	per              gatt.Peripheral
	colourChar       *gatt.Characteristic
	effectChar       *gatt.Characteristic
	toggleChar       *gatt.Characteristic
	doneChannel      chan struct{}
	connectedChannel chan bool
	connected        bool
}

func NewCandle(id string) *candle {
	return &candle{
		id:               id,
		doneChannel:      make(chan struct{}),
		connectedChannel: make(chan bool),
		connected:        false,
	}
}

func (p *candle) Off() error {
	c := NewColour(0, 0, 0, 0)
	e := p.solidColourEffect(c)

	return p.SetEffect(e)
}

func (p *candle) SetColour(c *Colour) error {
	e := p.solidColourEffect(c)
	payload := p.effectPayload(e)

	err := p.per.WriteCharacteristic(p.effectChar, payload, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *candle) ReadColour() ([]byte, error) {
	currentColour, err := p.per.ReadCharacteristic(p.colourChar)

	if err != nil {
		fmt.Println("Failed to read characteristic:", err)
		return nil, err
	}

	return currentColour, nil
}

func (p *candle) SetEffect(e *Effect) error {
	payload := p.effectPayload(e)

	err := p.per.WriteCharacteristic(p.effectChar, payload, true)
	if err != nil {
		return err
	}
	return nil
}

func (p *candle) ReadEffect() ([]byte, error) {
	currentEffect, err := p.per.ReadCharacteristic(p.effectChar)
	if err != nil {
		fmt.Println("Failed to read characteristic:", err)
		return nil, err
	}
	return currentEffect, nil
}

func (p *candle) Connect() error {
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

func (p *candle) OnToggle(f func(bool)) error {
	return p.per.SetNotifyValue(p.toggleChar, func(c *gatt.Characteristic, b []byte, err error) {
		if len(b) >= 4 {
			isOn := b[1] != 0 || b[2] != 0 || b[3] != 0
			f(isOn)
		}
	})
}

func (p *candle) IsConnected() bool {
	return p.connected
}

func (p *candle) Disconnect() {
	if p.connected {
		p.per.Device().CancelConnection(p.per)
	}
}

func (p *candle) colourPayload(c *Colour) []byte {
	return []byte{c.Brightness(), c.R(), c.G(), c.B()}
}

func (p *candle) effectPayload(e *Effect) []byte {
	return []byte{
		e.Colour().Brightness(), e.Colour().R(), e.Colour().G(), e.Colour().B(),
		e.Mode(), 0, e.Speed(), 0}
}

func (p *candle) solidColourEffect(c *Colour) *Effect {
	return NewEffect(SOLID, c, 0)
}

func (p *candle) onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func (p *candle) onPeripheralDiscovered(per gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if strings.ToUpper(per.ID()) != strings.ToUpper(p.id) {
		return
	}

	per.Device().StopScanning()
	per.Device().Connect(per)
}

func (p *candle) onPeripheralConnected(per gatt.Peripheral, err error) {
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
			case cUUID == toggleCharUUID:
				p.toggleChar = c
			default:
				continue
			}
		}
	}

	p.per = per

	p.connectedChannel <- true
}

func (p *candle) onPeripheralDisconnected(per gatt.Peripheral, err error) {
	p.connected = false

	p.per = nil
	p.colourChar = nil
	p.effectChar = nil
}

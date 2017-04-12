package playbulb

import (
  "strings"
  "fmt"
  "github.com/currantlabs/gatt"
  "github.com/currantlabs/gatt/examples/option"
)

const (
  // TODO turn these into UUIDs
  candleServiceUUID = "ff02"
  effectCharUUID = "fffb"
  colourCharUUID = "fffc"
)

// TODO add Logging with controllable level so the fmt.Print's can be removed
type PlaybulbCandle struct {
  id string
  per gatt.Peripheral
  colourChar  *gatt.Characteristic
  effectChar  *gatt.Characteristic
  doneChannel chan struct{}
  connectedChannel chan bool
  connected bool
}

func NewPlaybulbCandle(id string) *PlaybulbCandle {
  p := PlaybulbCandle {
    id: id,
    doneChannel: make(chan struct{}),
    connectedChannel: make(chan bool),
    connected: false,
  }
  return &p
}

func (p *PlaybulbCandle) Off() {
  c := NewColour(0, 0, 0, 0)
  e := p.solidColourEffect(c)

  p.SetEffect(e)
}

func (p *PlaybulbCandle) SetColour(c *Colour) error {
  e := p.solidColourEffect(c)
  payload := p.effectPayload(e)

  err := p.per.WriteCharacteristic(p.effectChar, payload, true)
  if err != nil {
    return err
  }
  return nil
}

func (p *PlaybulbCandle) ReadColour() ([]byte, error) {
  currentColour, err := p.per.ReadCharacteristic(p.colourChar)

  if err != nil {
    fmt.Println("Failed to read characteristic:", err)
    return nil, err
  }

  return currentColour, nil
}

func (p *PlaybulbCandle) SetEffect(e *Effect) error {
  payload := p.effectPayload(e)

  err := p.per.WriteCharacteristic(p.effectChar, payload, true)
  if err != nil {
    return err
  }
  return nil
}

func (p *PlaybulbCandle) ReadEffect() ([]byte, error) {
  currentEffect, err := p.per.ReadCharacteristic(p.effectChar)
  if err != nil {
    fmt.Println("Failed to read characteristic:", err)
    return nil, err
  }
  return currentEffect, nil
}

func (p *PlaybulbCandle) Connect() error {
  // TODO add a timeout if device not found within x seconds
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

  p.connected = <- p.connectedChannel
  return nil
}

func (p *PlaybulbCandle) Disconnect() {
  if p.connected {
    p.per.Device().CancelConnection(p.per)
  }
}

func (p *PlaybulbCandle) colourPayload(c *Colour) []byte {
  return []byte{ c.Brightness(), c.R(), c.G(), c.B() }
}

func (p *PlaybulbCandle) effectPayload(e *Effect) []byte {
  return []byte{
    e.Colour().Brightness(), e.Colour().R(), e.Colour().G(), e.Colour().B(),
    e.Mode(), 0, e.Speed(), 0 }
}

func (p *PlaybulbCandle) solidColourEffect(c *Colour) *Effect {
  return NewEffect(SOLID, c, 0)
}

func (p *PlaybulbCandle) onStateChanged(d gatt.Device, s gatt.State) {
  fmt.Println("State: ", s)
  switch s {
  case gatt.StatePoweredOn:
    fmt.Println("Scanning...")
    d.Scan([]gatt.UUID{}, false)
    return
  default:
    d.StopScanning()
  }
}

func (p *PlaybulbCandle) onPeripheralDiscovered(per gatt.Peripheral, a *gatt.Advertisement, rssi int) {
  if strings.ToUpper(per.ID()) != strings.ToUpper(p.id) {
    return
  }

  per.Device().StopScanning()
  per.Device().Connect(per)
}

func (p *PlaybulbCandle) onPeripheralConnected(per gatt.Peripheral, err error) {
  fmt.Println("Connected to playbulb", per.ID())

  services, err := per.DiscoverServices(nil) // p.DiscoverServices([]gatt.UUID{candleServiceUUID})
  if err != nil {
    fmt.Printf("Failed to discover services, err: %s\n", err)
    return
  }

  for _, s := range services {
    if s.UUID().String() != candleServiceUUID {
      continue
    }

    cs, err := per.DiscoverCharacteristics(nil, s) // per.DiscoverCharacteristics([]gatt.UUID{charEffectUUID,colourEffectUUID}, s)
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

func (p *PlaybulbCandle) onPeripheralDisconnected(per gatt.Peripheral, err error) {
  p.connectedChannel <- false

  p.per = nil
  p.colourChar = nil
  p.effectChar = nil
}
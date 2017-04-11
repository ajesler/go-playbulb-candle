
package main

import (
  "flag"
  "fmt"
  "log"
  "strings"
  "time"

  "github.com/currantlabs/gatt"
  "github.com/currantlabs/gatt/examples/option"
)

const (
  effect_candle    string = "candle"
  effect_solid     string = "solid"
  effect_rainbow   string = "rainbow"
)

var (
  bulbID = flag.String("id", "", "the id of the bulb (required)")
  effect = flag.String("effect", "", "[solid|candle|flash|pulse|rainbow|fade]")
  colour = flag.String("colour", "", "[red|green|blue]")

  done = make(chan struct{})
)

func prepareColourPayload() []byte {
  switch *colour {
  case "red":
    return []byte{'\x00', '\xff', '\x00', '\x00'}
  case "green":
    return []byte{'\x00', '\x00', '\xff', '\x00'}
  case "blue":
    return []byte{'\x00', '\x00', '\x00', '\xff'}
  default:
    return []byte{'\xff', '\x00', '\x00', '\x00'}
  }
}

func prepareEffectPayload() []byte {
  fmt.Println("Setting effect to", *effect)
  switch *effect {
  case effect_rainbow:
    fmt.Println("Effect set to", effect_rainbow)
    return []byte{'\x00', '\xff', '\xff', '\x00', '\x02', '\x00', '\x14', '\x00'}
  case effect_solid:
    fmt.Println("Effect set to", effect_solid)
    // 00ff0000 ff000100
    return []byte{'\x00', '\x00', '\x00', '\xff', '\xff', '\x00', '\x01', '\x00'}
  case effect_candle:
    fmt.Println("Effect set to", effect_candle)
    return []byte{'\x00', '\xff', '\x00', '\x00', '\x04', '\x00', '\x01', '\x00'}
  default:
    fmt.Println("Effect set to", effect_solid)
    return []byte{'\x00', '\x00', '\x00', '\x00', '\xff', '\x00', '\x40', '\x00'}
  }
}

func onStateChanged(device gatt.Device, state gatt.State) {
  fmt.Println("State: ", state)
  switch state {
  case gatt.StatePoweredOn:
    fmt.Println("Scanning...")
    device.Scan([]gatt.UUID{}, false)
    return
  default:
    device.StopScanning()
  }
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
  id := strings.ToUpper(*bulbID)
  if strings.ToUpper(p.ID()) != id {
    return
  }

  fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
  fmt.Println("  Local Name        =", a.LocalName)
  fmt.Println("  TX Power Level    =", a.TxPowerLevel)
  fmt.Println("  Manufacturer Data =", a.ManufacturerData)
  fmt.Println("  Service Data      =", a.ServiceData)

  p.Device().StopScanning()
  p.Device().Connect(p)
}

func onPeripheralConnected(p gatt.Peripheral, err error) {
  fmt.Println("Connected to playbulb", p.ID())

  // Discover services
  ss, err := p.DiscoverServices(nil)
  if err != nil {
    fmt.Printf("Failed to discover services, err: %s\n", err)
    return
  }

  deviceCharacteristics := make(map[string]*gatt.Characteristic)
  for _, s := range ss {
    fmt.Println("Found service", s.UUID().String())
    if s.UUID().String() != "ff02" {
      continue
    }

    cs, err := p.DiscoverCharacteristics(nil, s)
    if err != nil {
      fmt.Printf("Failed to discover characteristics, err: %s\n", err)
      continue
    }

    for _, c := range cs {
      characteristicName := "unknown"
      switch c.UUID().String() {
      case "fffb":
        characteristicName = "EFFECT"
      case "fffc":
        characteristicName = "COLOUR"
      default:
        continue
      }
      fmt.Println(characteristicName, "==>>", c.UUID().String())
      fmt.Printf("   %+v\n", c)
      deviceCharacteristics[characteristicName] = c
    }
  }

  if *colour != "" {
    colourPayload := prepareColourPayload()
    err = p.WriteCharacteristic(deviceCharacteristics["COLOUR"], colourPayload, true)
    if err != nil {
      fmt.Println("Failed to set colour:", err)
      return
    }

    time.Sleep(2 * time.Second)

    currentColour, err := p.ReadCharacteristic(deviceCharacteristics["COLOUR"])
    if err != nil {
      fmt.Println("Failed to read characteristic:", err)
      return
    }
    fmt.Printf("got colour: %x\n", currentColour)
  } else if *effect != "" {
    effectPayload := prepareEffectPayload()
    err = p.WriteCharacteristic(deviceCharacteristics["EFFECT"], effectPayload, true)
    if err != nil {
      fmt.Println("Failed to set effect:", err)
      return
    }

    time.Sleep(2 * time.Second)

    currentEffect, err := p.ReadCharacteristic(deviceCharacteristics["EFFECT"])
    if err != nil {
      fmt.Println("Failed to read characteristic:", err)
      return
    }
    fmt.Printf("got effect: %x\n", currentEffect)
  }

  time.Sleep(1 * time.Second)

  p.Device().CancelConnection(p)
}

func onPeripheralDisconnected(p gatt.Peripheral, err error) {
  close(done)
}

func main() {
  flag.Parse()

  if (*colour == "" && *effect == "") {
    flag.PrintDefaults()
    return
  }

  device, err := gatt.NewDevice(option.DefaultClientOptions...)
  if err != nil {
    log.Fatalf("Failed to open device, err: %s\n", err)
  }

  device.Handle(
    gatt.PeripheralDiscovered(onPeripheralDiscovered),
    gatt.PeripheralConnected(onPeripheralConnected),
    gatt.PeripheralDisconnected(onPeripheralDisconnected),
  )

  device.Init(onStateChanged)

  <-done
  fmt.Println("Complete")
}
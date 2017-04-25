package main

import (
	"flag"
	"fmt"
	"github.com/ajesler/playbulb-candle"
	"time"
)

const (
	candleID          = "7d61a760051544f4b83307bc316fada6" // Change this to your candle id
	effectChangeDelay = 4 * time.Second
)

func delay() {
	fmt.Println("Sleeping for", effectChangeDelay)
	time.Sleep(effectChangeDelay)
}

func main() {
	flag.Parse()

	candleID := flag.Args()
	if len(candleID) != 1 {
		fmt.Println("Please supply a candleID, eg 'go run examples/demo.go <your-candle-id>'")
		return
	}

	p := playbulb.NewCandle(candleID[0])

	fmt.Println("Connecting to", candleID)
	err := p.Connect()
	if err != nil {
		fmt.Println("Failed to connect", err)
	}

	p.OnToggle(func(isOn bool) {
		if isOn {
			fmt.Println("LED on")
		} else {
			fmt.Println("LED off")
		}
	})

	red, _ := playbulb.ColourFromHexString("00FF0000")
	green, _ := playbulb.ColourFromHexString("8000FF00") // light green
	blue, _ := playbulb.ColourFromHexString("000000FF")

	fmt.Println("-> Solid green")
	p.SetColour(green)

	delay()

	fmt.Println("-> Solid red")
	p.SetColour(red)

	delay()

	fmt.Println("-> Blue candle")
	e := playbulb.NewEffect(playbulb.CANDLE, blue, 1)
	p.SetEffect(e)

	delay()

	fmt.Println("-> Fast Rainbow")
	e = playbulb.NewEffect(playbulb.RAINBOW, red, 200)
	p.SetEffect(e)

	delay()

	fmt.Println("-> Blue pulse")
	e = playbulb.NewEffect(playbulb.PULSE, blue, 230)
	p.SetEffect(e)

	delay()

	effect, _ := p.ReadEffect()
	colour, _ := p.ReadColour()

	fmt.Printf("Current colour is %x\n", colour)
	fmt.Printf("Current effect is %x\n", effect)

	delay()

	fmt.Println("-> Off")
	p.Off()
	time.Sleep(time.Second) // this is needed to give the command time to send before we disconnect

	fmt.Println("Disconnecting")
	p.Disconnect()
	fmt.Println("Done")
}

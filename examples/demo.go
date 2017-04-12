package main

import (
	"fmt"
	"go-playbulb-candle"
	"time"
)

const (
	candleID          = "700ef5636bb249b8b95cdee4def26000" // Change this to your candle id
	effectChangeDelay = 4 * time.Second
)

func delay() {
	fmt.Println("Sleeping for", effectChangeDelay)
	time.Sleep(effectChangeDelay)
}

func main() {
	p := playbulb.NewCandle(candleID)

	fmt.Println("Connecting to", candleID)
	err := p.Connect()
	if err != nil {
		fmt.Println("Failed to connect", err)
	}

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

package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"go-playbulb-candle"
)

var (
	idFlag     = flag.String("id", "", "the id of the bulb (required)")
	effectFlag = flag.String("effect", "", "[flash|pulse|rainbow|fade|candle|solid]")
	colourFlag = flag.String("colour", "", "6 or 8 character hex code. If 8 characters, the first byte is the brightness")
	speedFlag  = flag.Int("speed", 0, "a value from 0 - 255")
)

func effectMode() (byte, error) {
	switch *effectFlag {
	case "flash":
		return playbulb.FLASH, nil
	case "pulse":
		return playbulb.PULSE, nil
	case "rainbow":
		return playbulb.RAINBOW, nil
	case "fade":
		return playbulb.FADE, nil
	case "candle":
		return playbulb.CANDLE, nil
	case "solid":
		return playbulb.SOLID, nil
	case "":
		return playbulb.SOLID, nil
	default:
		return 0, errors.New("Unsupported effect")
	}
}

func main() {
	flag.Parse()

	if *colourFlag == "" && *effectFlag == "" {
		flag.PrintDefaults()
		return
	}

	if *colourFlag == "" {
		*colourFlag = "00FF0000"
	}

	colour, err := playbulb.ColourFromHexString(*colourFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	eM, err := effectMode()
	if err != nil {
		fmt.Println(err)
		return
	}

	speed := byte(0)
	if *speedFlag > 255 || *speedFlag < 0 {
		fmt.Println("Speed must be between 0 and 255")
		return
	} else {
		speed = byte(*speedFlag)
	}

	effect := playbulb.NewEffect(eM, colour, speed)

	candle := playbulb.NewCandle(*idFlag)

	candle.Connect()
	candle.SetEffect(effect)

	// Required to give the SetEffect time to send before disconnection
	time.Sleep(time.Second)

	candle.Disconnect()
}

package playbulb

import (
	"encoding/hex"
	"errors"
	"regexp"
)

type Colour struct {
	brightness, r, g, b uint8
}

func NewColour(brightness uint8, r uint8, g uint8, b uint8) *Colour {
	c := Colour{brightness: brightness, r: r, g: g, b: b}
	return &c
}

func (c *Colour) Brightness() uint8 {
	return c.brightness
}

func (c *Colour) R() uint8 {
	return c.r
}

func (c *Colour) G() uint8 {
	return c.g
}

func (c *Colour) B() uint8 {
	return c.b
}

func ColourFromHexString(s string) (*Colour, error) {
	var br, r, g, b uint8

	validColourString, _ := regexp.MatchString("^[a-zA-Z0-9]{6}([a-zA-Z0-9]{2})?$", s)
	if !validColourString {
		return nil, errors.New("Only 6 or 8 character hex colours are supported")
	}

	if len(s) == 8 {
		br = hexToUint8(s[:2])
		s = s[2:]
	}

	r = hexToUint8(s[:2])
	g = hexToUint8(s[2:4])
	b = hexToUint8(s[4:])

	return NewColour(br, r, g, b), nil
}

func hexToUint8(s string) uint8 {
	v, _ := hex.DecodeString(s)
	return v[0]
}

package playbulb

import (
	"encoding/hex"
	"errors"
	"regexp"
)

type Colour struct {
	brightness, r, g, b byte
}

func NewColour(brightness byte, r byte, g byte, b byte) *Colour {
	c := Colour{brightness: brightness, r: r, g: g, b: b}
	return &c
}

func (c *Colour) Brightness() byte {
	return c.brightness
}

func (c *Colour) R() byte {
	return c.r
}

func (c *Colour) G() byte {
	return c.g
}

func (c *Colour) B() byte {
	return c.b
}

func ColourFromHexString(s string) (*Colour, error) {
	var br, r, g, b byte

	validColourString, _ := regexp.MatchString("^[a-zA-Z0-9]{6}([a-zA-Z0-9]{2})?$", s)
	if !validColourString {
		return nil, errors.New("Only 6 or 8 character hex colours are supported")
	}

	if len(s) == 8 {
		br = hexToByte(s[:2])
		s = s[2:]
	}

	r = hexToByte(s[:2])
	g = hexToByte(s[2:4])
	b = hexToByte(s[4:])

	return NewColour(br, r, g, b), nil
}

func hexToByte(s string) byte {
	v, _ := hex.DecodeString(s)
	return v[0]
}

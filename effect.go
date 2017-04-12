package playbulb

const (
	FLASH   byte = 0
	PULSE   byte = 1
	RAINBOW byte = 2
	FADE    byte = 3
	CANDLE  byte = 4
	SOLID   byte = 5
)

type Effect struct {
	mode   byte
	colour *Colour
	speed  byte
}

func NewEffect(m byte, c *Colour, s byte) *Effect {
	e := Effect{mode: m, colour: c, speed: s}
	return &e
}

func (e *Effect) Mode() byte {
	return e.mode
}

func (e *Effect) Colour() *Colour {
	return e.colour
}

func (e *Effect) Speed() byte {
	return e.speed
}

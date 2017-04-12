package playbulb

const (
	FLASH   uint8 = 0
	PULSE   uint8 = 1
	RAINBOW uint8 = 2
	FADE    uint8 = 3
	CANDLE  uint8 = 4
	SOLID   uint8 = 5
)

type Effect struct {
	mode   uint8
	colour *Colour
	speed  uint8
}

func NewEffect(m uint8, c *Colour, s uint8) *Effect {
	e := Effect{mode: m, colour: c, speed: s}
	return &e
}

func (e *Effect) Mode() uint8 {
	return e.mode
}

func (e *Effect) Colour() *Colour {
	return e.colour
}

func (e *Effect) Speed() uint8 {
	return e.speed
}

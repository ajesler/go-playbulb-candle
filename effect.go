package playbulb

type Mode int

const (
  CANDLE  Mode = 1
  FLASH   Mode = 2
  PULSE   Mode = 3
  RAINBOW Mode = 4
  FADE    Mode = 5
)

type Effect struct {
  mode Mode
  colour *Colour
  speed uint8
}

func NewEffect(m Mode, c *Colour, s uint8) *Effect {
  e := Effect { mode: m, colour: c, speed: s }
  return &e
}

func (e *Effect) Mode() Mode {
  return e.mode
}

func (e *Effect) Colour() *Colour {
  return e.colour
}

func (e *Effect) Speed() uint8 {
  return e.speed
}
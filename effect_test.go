package playbulb_test

import (
	"github.com/ajesler/playbulb-candle"
	"testing"
)

func TestEffect_NewEffect(t *testing.T) {
	m := playbulb.PULSE
	c := playbulb.NewColour(1, 2, 3, 4)
	s := byte(38)

	e := playbulb.NewEffect(m, c, s)

	if e.Mode() != playbulb.PULSE || e.Colour() != c || e.Speed() != s {
		t.Error("expected effect values did not match actual. m:", e.Mode(), "c: ", e.Colour(), "s:", e.Speed())
	}
}

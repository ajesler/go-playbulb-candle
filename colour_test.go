package playbulb_test

import (
  "go-playbulb-candle"
  "testing"
)

func TestColour_NewColour(t *testing.T) {
  c := playbulb.NewColour(0, 254, 0, 25)

  if c.Brightness() != 0 || c.R() != 254 || c.G() != 0 || c.B() != 25 {
    t.Error("expected colour fields did not match actual. brightness:", c.Brightness(), "r:", c.R(), "g:",c.G(), "b:", c.B())
  }
}

func TestColour_FromString_8length(t *testing.T) {
  c, _ := playbulb.FromString("08AC00D1")

  if c.Brightness() != 8 || c.R() != 172 || c.G() != 0 || c.B() != 209 {
    t.Error("expected colour fields did not match actual. brightness:", c.Brightness(), "r:", c.R(), "g:",c.G(), "b:", c.B())
  }
}

func TestColour_FromString_6length(t *testing.T) {
  c, _ := playbulb.FromString("AC00D1")

  if c.Brightness() != 0 || c.R() != 172 || c.G() != 0 || c.B() != 209 {
    t.Error("expected colour fields did not match actual. brightness:", c.Brightness(), "r:", c.R(), "g:",c.G(), "b:", c.B())
  }
}

func TestColour_FromString_error(t *testing.T) {
  _, err := playbulb.FromString("AC00D10")

  if err == nil {
    t.Error("Expected an error, but didn't get one")
  }
}
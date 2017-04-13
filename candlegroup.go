package playbulb

import (
  "errors"
  "strings"
)

type CandleGroup struct {
  candles []Candle
}

func NewCandleGroup(cs ...Candle) *CandleGroup {
  return &CandleGroup{
    candles: cs,
  }
}

func (cg *CandleGroup) Connect() error {
  es := make([]string, 0)

  for _, c := range cg.candles {
    e := c.Connect()
    if e != nil {
      es = append(es, e.Error())
    }
  }

  if len(es) > 0 {
    return errors.New(strings.Join(es, "\n"))
  } else {
    return nil
  }
}

func (cg *CandleGroup) IsConnected() bool {
  for _, c := range cg.candles {
    if c.IsConnected() {
        return true
    }
  }
  return false
}

func (cg *CandleGroup) Disconnect() {
  for _, c := range cg.candles {
    c.Disconnect()
  }
}

func (cg *CandleGroup) SetEffect(e *Effect) error {
  es := make([]string, 0)

  for _, c := range cg.candles {
    e := c.SetEffect(e)
    if e != nil {
      es = append(es, e.Error())
    }
  }

  if len(es) > 0 {
    return errors.New(strings.Join(es, "\n"))
  } else {
    return nil
  }
}

func (cg *CandleGroup) SetColour(cl *Colour) error {
  es := make([]string, 0)

  for _, c := range cg.candles {
    e := c.SetColour(cl)
    if e != nil {
      es = append(es, e.Error())
    }
  }

  if len(es) > 0 {
    return errors.New(strings.Join(es, "\n"))
  } else {
    return nil
  }
}

func (cg *CandleGroup) Off() error {
  es := make([]string, 0)

  for _, c := range cg.candles {
    e := c.Off()
    if e != nil {
      es = append(es, e.Error())
    }
  }

  if len(es) > 0 {
    return errors.New(strings.Join(es, "\n"))
  } else {
    return nil
  }
}


// Code generated by "go generate github.com/jak9708/gonummat/unit”; DO NOT EDIT.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

// Capacitance represents an electrical capacitance in Farads.
type Capacitance float64

const Farad Capacitance = 1

// Unit converts the Capacitance to a *Unit.
func (cp Capacitance) Unit() *Unit {
	return New(float64(cp), Dimensions{
		CurrentDim: 2,
		LengthDim:  -2,
		MassDim:    -1,
		TimeDim:    4,
	})
}

// Capacitance allows Capacitance to implement a Capacitancer interface.
func (cp Capacitance) Capacitance() Capacitance {
	return cp
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (cp *Capacitance) From(u Uniter) error {
	if !DimensionsMatch(u, Farad) {
		*cp = Capacitance(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*cp = Capacitance(u.Unit().Value())
	return nil
}

func (cp Capacitance) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", cp, float64(cp))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " F"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(cp))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(cp))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(cp))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(cp))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g F)", c, cp, float64(cp))
	}
}

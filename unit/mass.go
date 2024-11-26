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

// Mass represents a mass in kilograms.
type Mass float64

const (
	Gram Mass = 1e-3

	Kilogram = Kilo * Gram
)

// Unit converts the Mass to a *Unit.
func (m Mass) Unit() *Unit {
	return New(float64(m), Dimensions{
		MassDim: 1,
	})
}

// Mass allows Mass to implement a Masser interface.
func (m Mass) Mass() Mass {
	return m
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (m *Mass) From(u Uniter) error {
	if !DimensionsMatch(u, Gram) {
		*m = Mass(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*m = Mass(u.Unit().Value())
	return nil
}

func (m Mass) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", m, float64(m))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " kg"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(m))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(m))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(m))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(m))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g kg)", c, m, float64(m))
	}
}

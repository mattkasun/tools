//nolint:mnd
package money

import (
	"fmt"
	"math"
)

// money represent a monetary value in decimal format, ie $1 = 100.
type money int64

// String implements the stringer interface for Money.
func (m money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	dollars := m / 100
	cents := m % 100
	dStrg := fmt.Sprintf("%d", dollars)
	for i := len(dStrg) - 3; i > 0; i -= 3 {
		dStrg = dStrg[:i] + "," + dStrg[i:]
	}
	return fmt.Sprintf("$%s%s.%02d", sign, dStrg, cents)
}

// Tax calculates the amount of tax given the rate on a Money amount.
func (m money) Tax(rate float64) money {
	return money(math.Round(float64(m) * rate))
}

// WithTax returns the amount with tax included.
func (m money) WithTax(rate float64) money {
	tax := m.Tax(rate)
	return m + tax
}

// New returns a Money representation of a float; max value is one hundred trillon.
func New(amount float64) money { //nolint:revive
	if amount >= 1e16 {
		return money(1e16)
	}
	if amount <= -1e16 {
		return money(-1e16)
	}
	return money(math.Round(amount * 100))
}

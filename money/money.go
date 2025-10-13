// Package money is a simple currency package
package money

import (
	"fmt"
	"math"
)

const (
	maxValue      = 1e16
	dollarInCents = 100
	three         = 3
)

// Money represent a monetary value in decimal format, ie $1 = 100.
type Money int64

// String implements the stringer interface for Money.
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	dollars := m / dollarInCents
	cents := m % dollarInCents
	dollar := fmt.Sprintf("%d", dollars)
	for i := len(dollar) - three; i > 0; i -= 3 {
		dollar = dollar[:i] + "," + dollar[i:]
	}
	return fmt.Sprintf("$%s%s.%02d", sign, dollar, cents)
}

// Tax calculates the amount of tax given the rate on a Money amount.
func (m Money) Tax(rate float64) Money {
	return Money(math.Round(float64(m) * rate))
}

// WithTax returns the amount with tax included.
func (m Money) WithTax(rate float64) Money {
	tax := m.Tax(rate)
	return m + tax
}

// New returns a Money representation of a float; max value is one hundred trillon.
func New(amount float64) Money {
	if amount >= maxValue {
		return Money(maxValue)
	}
	if amount <= -maxValue {
		return Money(-maxValue)
	}
	return Money(math.Round(amount * dollarInCents))
}

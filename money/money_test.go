package money_test

import (
	"math"
	"testing"

	"github.com/Kairum-Labs/should"
	"github.com/mattkasun/tools/money"
)

func TestString(t *testing.T) {
	should.BeEqual(t, money.Money(0).String(), "$0.00")
	should.BeEqual(t, money.Money(123456).String(), "$1,234.56")
	should.BeEqual(t, money.Money(-789).String(), "$-7.89")
	should.BeEqual(t, money.Money(math.MaxInt64).String(), "$92,233,720,368,547,758.07")
}

func TestTax(t *testing.T) {
	should.BeEqual(t, money.Money(10000).Tax(0.15), money.Money(1500))
	should.BeEqual(t, money.Money(5000).Tax(0), money.Money(0))
	should.BeEqual(t, money.Money(-10000).Tax(0.10), money.Money(-1000))
	should.BeEqual(t, money.Money(999).Tax(0.05), money.Money(50)) // 49.95 → rounds to 50
}

func TestWithTax(t *testing.T) {
	should.BeEqual(t, money.New(10000).WithTax(0.15), money.New(11500))
	should.BeEqual(t, money.New(20000).WithTax(0), money.New(20000))
	should.BeEqual(t, money.New(-10000).WithTax(0.10), money.New(-11000))
}

func TestNew(t *testing.T) {
	should.BeEqual(t, money.New(10.00), money.Money(1000))
	should.BeEqual(t, money.New(10.125), money.Money(1013)) // rounds correctly
	overflow := float64(math.MaxInt64) + 1
	should.BeEqual(t, money.New(overflow), money.New(math.MaxInt64))
	underflow := float64(math.MinInt64)
	should.BeEqual(t, money.New(underflow), money.Money(-1e16))
}

func TestRoundingBehavior(t *testing.T) {
	should.BeEqual(t, money.New(0.005), money.Money(1)) // half-up rounding
	should.BeEqual(t, money.New(0.004), money.Money(0))
}

package money_test

import (
	"math"
	"testing"

	"github.com/Kairum-Labs/should"
	"github.com/mattkasun/tools/money"
)

func TestString(t *testing.T) {
	should.BeEqual(t, money.New(0).String(), "$0.00")
	should.BeEqual(t, money.New(1234.56).String(), "$1,234.56")
	should.BeEqual(t, money.New(-7.89).String(), "$-7.89")
	should.BeEqual(t, money.New(math.MaxInt64).String(), "$100,000,000,000,000.00")
}

func TestTax(t *testing.T) {
	should.BeEqual(t, money.New(10000).Tax(0.15), money.New(1500))
	should.BeEqual(t, money.New(5000).Tax(0), money.New(0))
	should.BeEqual(t, money.New(-10000).Tax(0.10), money.New(-1000))
	should.BeEqual(t, money.New(9.99).Tax(0.05), money.New(.50)) // 49.95 → rounds to 50
}

func TestWithTax(t *testing.T) {
	should.BeEqual(t, money.New(10000).WithTax(0.15), money.New(11500))
	should.BeEqual(t, money.New(20000).WithTax(0), money.New(20000))
	should.BeEqual(t, money.New(-10000).WithTax(0.10), money.New(-11000))
}

func TestNew(t *testing.T) {
	should.BeEqual(t, int64(money.New(10.00)), int64(1000))
	should.BeEqual(t, int64(money.New(10.125)), int64(money.New(10.13))) // rounds correctly
	overflow := float64(math.MaxInt64) + 1
	should.BeEqual(t, money.New(overflow), money.New(math.MaxInt64))
}

func TestRoundingBehavior(t *testing.T) {
	should.BeEqual(t, money.New(0.005), money.New(0.01)) // half-up rounding
	should.BeEqual(t, money.New(0.004), money.New(0))
}

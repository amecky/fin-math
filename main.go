package main

import (
	"fmt"

	"github.com/amecky/fin-math/math"
)

func main() {
	option := &math.Option{
		StrikePrice:      90.0,
		TimeToExpiration: 45,
		Type:             "PUT",
	}

	underlying := &math.Underlying{
		Symbol:     "Delivery Hero",
		Price:      85.52,
		Volatility: .4,
	}

	bs := math.NewBlackScholes(option, underlying, .0102)

	fmt.Println("delta", bs.Delta)
	fmt.Println("IV", bs.ImpliedVolatility)
	fmt.Println("Theo price", bs.TheoPrice)
	fmt.Println("Theta", bs.Theta)
}

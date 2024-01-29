package main

import (
	"fmt"

	"github.com/ch1s7ya/hermes/internal/currencies"
)

func main() {
	fmt.Println("Hello, Bender!")

	// exchange_rate := currencies.GetExRate()
	// fmt.Println(string(exchange_rate))
	exchangeRate := currencies.GetExRate()
	fmt.Printf("Exchange rate:\n%v", exchangeRate)
}

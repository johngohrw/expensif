package domain

import "strconv"

var currencySymbols = map[string]string{
	"USD": "$",
	"MYR": "RM",
	"JPY": "¥",
	"CNY": "¥",
	"THB": "฿",
	"EUR": "€",
	"GBP": "£",
	"SGD": "S$",
	"KRW": "₩",
	"AUD": "A$",
	"CAD": "C$",
	"INR": "₹",
	"VND": "₫",
	"PHP": "₱",
	"IDR": "Rp",
	"HKD": "HK$",
	"TWD": "NT$",
}

func CurrencySymbol(code string) string {
	if s, ok := currencySymbols[code]; ok {
		return s
	}
	return "$"
}

// currencyDecimals lists ISO 4217 exceptions to the default 2 minor-unit
// decimal places. JPY, KRW, and VND have no minor unit in everyday use.
var currencyDecimals = map[string]int{
	"JPY": 0,
	"KRW": 0,
	"VND": 0,
}

// CurrencyDecimals returns the number of minor-unit decimal places for a
// currency code (defaults to 2, the ISO 4217 majority case).
func CurrencyDecimals(code string) int {
	if d, ok := currencyDecimals[code]; ok {
		return d
	}
	return 2
}

// AllCurrencyDecimals returns the decimal-place count for every supported
// currency, keyed by code — used to drive client-side amount formatting.
func AllCurrencyDecimals() map[string]int {
	result := make(map[string]int, len(currencySymbols))
	for code := range currencySymbols {
		result[code] = CurrencyDecimals(code)
	}
	return result
}

// FormatAmount formats amount using the currency's minor-unit decimal places.
func FormatAmount(amount float64, code string) string {
	return strconv.FormatFloat(amount, 'f', CurrencyDecimals(code), 64)
}

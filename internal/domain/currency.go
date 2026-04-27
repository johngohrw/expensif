package domain

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

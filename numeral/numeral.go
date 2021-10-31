package numeral

import (
	"strings"
)

type RomanNumeral struct {
    Value  int
    Symbol string
}

var allRomanNumerals = RomanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
    {10, "X"},
    {9, "IX"},
    {5, "V"},
    {4, "IV"},
    {1, "I"},
}

func ConvertToRoman(arabic int) string {

	var result strings.Builder
    for _, numeral := range allRomanNumerals {
        for arabic >= numeral.Value {
            result.WriteString(numeral.Symbol)
            arabic -= numeral.Value
        }
    }
    // for arabic > 0 {
	// 	switch {
	// 	case arabic > 9:
	// 		result.WriteString("X")
	// 		arabic -= 10
	// 	case arabic > 8:
	// 		result.WriteString("IX")
	// 		arabic -= 9
	// 	case arabic > 4:
	// 		result.WriteString("V")
	// 		arabic -= 5
	// 	case arabic > 3:
	// 		result.WriteString("IV")
	// 		arabic -= 4
	// 	default:
	// 		result.WriteString("I")
	// 		arabic--
	// 	}
    // }
    return result.String()
}

type RomanNumerals []RomanNumeral

// When you index strings in Go, you get a byte. -> string([]byte{symbol})
func (r RomanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
    for _, s := range r {
        if s.Symbol == symbol {
            return s.Value
        }
    }
    return 0
}

func (r RomanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
}

type windowedRoman string

func (w windowedRoman) Symbols() (symbols [][]byte) {
    for i := 0; i < len(w); i++ {
        symbol := w[i]
        notAtEnd := i+1 < len(w)
        if notAtEnd && isSubtractive(symbol) && allRomanNumerals.Exists(symbol, w[i+1]) {
            symbols = append(symbols, []byte{symbol, w[i+1]})
            i++
        } else {
            symbols = append(symbols, []byte{symbol})
        }
    }
    return
}

func ConvertToArabic(roman string) (total int) {
	for _, symbols := range windowedRoman(roman).Symbols() {
		total += allRomanNumerals.ValueOf(symbols...)
	}
	return
}

func isSubtractive(symbol uint8) bool {
	return symbol == 'I' || symbol == 'X' || symbol =='C'
}
// https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/roman-numerals

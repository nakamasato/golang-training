package numeral

import (
	"strings"
)

type RomanNumeral struct{
	Value uint16
	Symbol string
}

type RomanNumerals []RomanNumeral

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

func ConvertToRoman(arabic uint16) string {

	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}
	return result.String()
}

func (r RomanNumerals) ValueOf(roman ...byte) uint16 {
	for _, s := range r {
		if s.Symbol == string(roman) {
			return s.Value
		}
	}
	return 0
}

func (r RomanNumerals) Exists(symbols ...byte) bool {
	for _, s := range r {
		if  s.Symbol == string(symbols) {
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

func ConvertToArabic(roman string) (total uint16) {
	// this cannot check if subtractive
	// for _, symbol := range roman {
	// 	total += allRomanNumerals.ValueOf(string(symbol))
	// 	fmt.Printf("%s", string(symbol))
	// }
	for _, symbols := range windowedRoman(roman).Symbols() {
        total += allRomanNumerals.ValueOf(symbols...)
    }
	return
}

func isSubtractive(symbol uint8) bool {
	return symbol == 'I' || symbol == 'X' || symbol == 'C'
}

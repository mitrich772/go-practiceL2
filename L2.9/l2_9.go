package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func unpack(s string) (string, error) {

	var builder strings.Builder
	var prevRune rune
	esc := false //Флаг экранирования

	for _, r := range s {
		switch {
		case esc: // Если предыдущий сивол был экранирован
			builder.WriteRune(r)
			prevRune = r
			esc = false
		case r == '\\': // Если экранирование
			esc = true
		case unicode.IsDigit(r): //Ежели цифра
			if prevRune == 0 {
				return "", errors.New("без символа но цифра есть")
			}

			builder.WriteString(strings.Repeat(string(prevRune), int(r-'0')-1))
		default: //Букова
			builder.WriteRune(r)
			prevRune = r
		}

	}
	if esc {
		return "", errors.New("некорректная строка: оканчивается на экранирование")
	}
	return builder.String(), nil
}

func main() {
	fmt.Println(unpack("45"))
}

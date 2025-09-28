package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Вход: "a4bc2d5e"
// Выход: "aaaabccddddde"
func unpack(s string) string {
	runes := []rune(s)
	var builder strings.Builder
	for i, v := range runes {
		if unicode.IsDigit(v) {
			if runes[i-1] == '/' {
				builder.WriteRune(v)
			}else{
				builder.WriteString(strings.Repeat(string(runes[i-1]), int(v-'0')-1))
			}
		}else {
			if v != '/'{
				builder.WriteRune(v)
			}
		}
	}
	return builder.String()
}

func main() {
	fmt.Println(unpack("qwe/4/5"))
}

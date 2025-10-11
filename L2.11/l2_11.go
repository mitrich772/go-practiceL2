package main

import (
	"fmt"
	"slices"
	"strings"
)

func findAnogramm(words []string) map[string][]string {
	res := make(map[string][]string)
	sortedWordToRealKey := make(map[string]string)

	for _, w := range words {
		word := strings.ToLower(w)
		sortedWord := getSotredString(word)
		realKey, ok := sortedWordToRealKey[sortedWord] //Видели ли такую группу

		if !ok { // Если впервые видим такую группу
			sortedWordToRealKey[sortedWord] = word // Добавляем сортрованное значение чтобы потом понимать что куда кидать если не первый раз появилась грпупа
		} else { //Если такая группа aнограмм уже есть
			if _, exists := res[realKey]; !exists { // Тк мы уже видели такую группу то пишем 1 слово которое равно ключу
				res[realKey] = append(res[realKey], realKey)
			}
			res[realKey] = append(res[realKey], word) // Добавляем текущее слово
		}
	}

	for key := range res { // Сортируем слова в каждой группе
		slices.Sort(res[key])
	}
	return res
}

func getSotredString(s string) string {
	runes := []rune(s)
	slices.SortFunc(runes, func(a, b rune) int {
		if a == b {
			return 0
		}
		if a > b {
			return 1
		}
		return -1
	})
	return string(runes)
}
func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Print(findAnogramm(words))
}

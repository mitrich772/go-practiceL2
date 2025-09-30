package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Читает все непустые строки с файла
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" { // пропускаем пустые строки
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Удаляет дубликаты но только в отсортированном массиве где дубликаты стоят рядом
func removeDuplicatesSorted(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	uniq := []string{lines[0]}
	prev := lines[0]

	for _, line := range lines[1:] {
		if line != prev {
			uniq = append(uniq, line)
			prev = line
		}
	}

	return uniq
}

func main() {
	var column int
	var numeric bool
	var reverse bool
	var unique bool

	flag.IntVar(&column, "k", 1, "column for sort")
	flag.BoolVar(&numeric, "n", false, "sort string as number")
	flag.BoolVar(&reverse, "r", false, "reverse sort")
	flag.BoolVar(&unique, "u", false, "only unique values")
	flag.Parse()

	fmt.Printf("column %d, numeric %t, reverse %t, unique %t\n", column, numeric, reverse, unique)

	lines, err := readLines("data.txt")
	if err != nil {
		log.Println(err)
	}

	slices.SortFunc(lines, func(a, b string) int { // правило сортировки
		fa := strings.Split(a, "\t")
		fb := strings.Split(b, "\t")

		if len(fa) < column || len(fb) < column {
			log.Printf("строка имеет меньше столбцов, чем k=%d: '%v' / '%v'\n", column, fa, fb)
			os.Exit(1)
		}

		va := fa[column-1]
		vb := fb[column-1]

		if numeric {
			aInt, _ := strconv.Atoi(va)
			bInt, _ := strconv.Atoi(vb)
			if reverse {
				if aInt < bInt {
					return 1
				} else if aInt > bInt {
					return -1
				}
				return 0
			} else {
				if aInt < bInt {
					return -1
				} else if aInt > bInt {
					return 1
				}
				return 0
			}
		} else {
			if reverse {
				return strings.Compare(vb, va)
			}
			return strings.Compare(va, vb)
		}
	})

	if unique { // уникальные значния
		lines = removeDuplicatesSorted(lines)
	}

	for _, v := range lines { //Вывод строк
		fmt.Println(v)
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Finder struct {
	afterMatch  int
	beforeMatch int
	ignoreReg   bool
}

func ReadLines(filename string) ([]string, error) {
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
func ParseFilenameExpr() (filename string, expr string) {
	filename = parseFilename()
	expr = parseExpr()
	os.Args = os.Args[:len(os.Args)-2]
	return
}
func parseExpr() (expr string) {
	argsLen := len(os.Args)
	expr = os.Args[argsLen-2]
	return
}
func parseFilename() (filename string) {
	argsLen := len(os.Args)
	filename = os.Args[argsLen-1]
	return
}

// Выводит массив строк
func PrintLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}

func (fn Finder) findWithContex(lines []string, re *regexp.Regexp) (result []string) {
	for i, line := range lines {
		if re.FindAllString(line, -1) != nil {
			if (fn.beforeMatch > i) && (fn.afterMatch > (len(lines)-1)-i) { // Если в начале и конце мало строк (да так тоже может быть)
				result = append(result, lines[:i+1]...)
				result = append(result, lines[i:]...)
			} else if fn.afterMatch > (len(lines)-1)-i { // Если в конце мало строк
				result = append(result, lines[i-fn.beforeMatch:i]...) // Добавляем перед совпадением не включая совпадение
				result = append(result, lines[i:]...)
			} else if fn.beforeMatch > i { // Если в начале мало строк
				result = append(result, lines[:i+1]...)
				result = append(result, lines[i:i+fn.afterMatch+1]...) // Добавляем после совпадения включая совпадение
			} else {
				result = append(result, lines[i-fn.beforeMatch:i]...)  // Добавляем перед совпадением не включая совпадение
				result = append(result, lines[i:i+fn.afterMatch+1]...) // Добавляем после совпадения включая совпадение
			}

			result = append(result, "-----------------")
		}
	}
	return
}
func (fn Finder) FindStrings(lines []string, template string) ([]string, error) {
	var result []string
	// Если надо игнорировать регистр
	if fn.ignoreReg {
		template = `(?i)` + template
	}
	// Готовим регулярку
	re, err := regexp.Compile(template)
	if err != nil {
		return nil, err
	}
	if fn.beforeMatch > 0 || fn.afterMatch > 0 { // До или после контекст
		result = append(result, fn.findWithContex(lines, re)...)
	} else { //Контекст до и после равен 0 меньше нуля не может быть (Выводим только строку с совпадением)
		re, err := regexp.Compile(template)
		if err != nil {
			return nil, err
		}
		for _, line := range lines {
			if re.FindAllString(line, -1) != nil {
				result = append(result, line)
			}
		}
	}
	return result, nil
}

func main() {
	var fn Finder
	var aroundMatch int
	filename, expression := ParseFilenameExpr()

	flag.BoolVar(&fn.ignoreReg, "i", false, "Ignore register")
	flag.IntVar(&fn.afterMatch, "A", 0, "Print N strings after match")
	if fn.afterMatch < 0 {
		log.Println("-A N должно получить положительное число, так что я исправлю на ноль, типо ничего не было ;)")
	}
	flag.IntVar(&fn.beforeMatch, "B", 0, "Print N strings before match")
	if fn.afterMatch < 0 {
		log.Println("-B N должно получить положительное число, так что я исправлю на ноль, типо ничего не было ;)")
	}
	flag.IntVar(&aroundMatch, "C", 0, "Print N strings around match")
	if fn.afterMatch < 0 {
		log.Println("-A N должно получить положительное число, так что я исправлю на ноль, типо ничего не было ;)")
	}
	flag.Parse()
	// Перебиваем -C в -B -A
	fn.afterMatch = aroundMatch
	fn.beforeMatch = aroundMatch

	fmt.Printf("Имя файла: %s\nВыражение для поиска: %s\n", filename, expression)

	lines, err := ReadLines(filename)
	if err != nil {
		log.Fatal(err)
	}

	findedLines, err := fn.FindStrings(lines, expression)
	if err != nil {
		log.Fatal(err)
	}
	PrintLines(findedLines)
}

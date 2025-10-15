package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Finder struct {
	afterMatch  int
	beforeMatch int
	ignoreReg   bool
	invert      bool
	addStrNum   bool
	onlyDigit   bool
	fixed       bool
}
type Line struct {
	Num  int
	Text string
}

func ReadLines(filename string) ([]Line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []Line
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" { // пропускаем пустые строки
			lines = append(lines, Line{
				Num:  lineNum,
				Text: line,
			})
		}
		lineNum++
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
func (fn Finder) PrintLines(lines []Line) {
	for _, line := range lines {
		if fn.addStrNum {
			fmt.Printf("%d ", line.Num)
		}
		fmt.Println(line.Text)
	}
}

// Только число совпадений ввиде Line
func (fn Finder) findOnlydigit(lines []Line, re *regexp.Regexp) (result Line) {
	digit := 0
	for _, line := range lines {
		if fn.matchWithInvertFlag(line, re) {
			digit++
		}
	}
	result = Line{Text: strconv.Itoa(digit)}
	return
}

// Ищет строки + контекст -A -B
func (fn Finder) findWithContex(lines []Line, re *regexp.Regexp) (result []Line) {
	for i, line := range lines {
		if fn.matchWithInvertFlag(line, re) { // Тут по inverted
			//log.Printf("%d", len(re.FindAllString(line, -1)))
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

			result = append(result, Line{Text: "-----------------"})
		}
	}
	return
}

// Ищет строки без контекста
func (fn Finder) findNoContex(lines []Line, re *regexp.Regexp) (result []Line) {
	for _, line := range lines {
		if fn.matchWithInvertFlag(line, re) {
			result = append(result, line)
		}
	}
	return
}

// Выдает bool для строк которые содержат/не содержат match от fn.inverted false/true
func (fn Finder) matchWithInvertFlag(line Line, re *regexp.Regexp) (res bool) {
	return re.MatchString(line.Text) != fn.invert
}

// FindStrings ищет строки по паттерну в зависимости от заданных параметров
func (fn Finder) FindStrings(lines []Line, template string) ([]Line, error) {
	var result []Line
	// Если надо воспринимать буквально
	if fn.fixed {
		template = regexp.QuoteMeta(template)
	}
	// Если надо игнорировать регистр
	if fn.ignoreReg {
		template = `(?i)` + template
	}
	// Готовим регулярку
	re, err := regexp.Compile(template)
	if err != nil {
		return nil, err
	}
	if fn.onlyDigit { // если надо только число вывести
		result = append(result, fn.findOnlydigit(lines, re))
	} else if fn.beforeMatch > 0 || fn.afterMatch > 0 { // До или после контекст
		result = append(result, fn.findWithContex(lines, re)...)
	} else { //Контекст до и после равен 0 меньше нуля не может быть (Выводим только строку с совпадением)
		result = append(result, fn.findNoContex(lines, re)...)
	}
	return result, nil
}

func main() {
	var fn Finder
	var aroundMatch int
	filename, expression := ParseFilenameExpr()

	flag.BoolVar(&fn.ignoreReg, "i", false, "Ignore register")
	flag.BoolVar(&fn.invert, "v", false, "invert, print strings with no match")
	flag.BoolVar(&fn.addStrNum, "n", false, "add string number in macth string start")
	flag.BoolVar(&fn.onlyDigit, "c", false, "print only number of matches")
	flag.BoolVar(&fn.fixed, "F", false, "seek full match")
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
	if aroundMatch > 0 {
		fn.afterMatch = aroundMatch
		fn.beforeMatch = aroundMatch
	}

	fmt.Printf("Имя файла: %s\nВыражение для поиска: %s\n", filename, expression)

	lines, err := ReadLines(filename)
	if err != nil {
		log.Fatal(err)
	}

	findedLines, err := fn.FindStrings(lines, expression)
	if err != nil {
		log.Fatal(err)
	}
	fn.PrintLines(findedLines)
}

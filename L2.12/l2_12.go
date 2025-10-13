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
	afterMatch int
	ignoreReg  bool
}

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
func parseFilenameExpr() (filename string, expr string) {
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
func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}
func (fn Finder) findStrings(lines []string, template string) ([]string, error) {
	var result []string
	// Если надо игнорировать регистр
	if fn.ignoreReg {
		template = `(?i)` + template
	}
	if fn.afterMatch > 0 {
		re, err := regexp.Compile(template)
		if err != nil {
			return nil, err
		}
		for i, line := range lines {
			if re.FindAllString(line, -1) != nil {
				if fn.afterMatch > i {
					result = append(result, lines[:i+1]...)
				} else {
					result = append(result, lines[i-fn.afterMatch:i+1]...)
				}

				result = append(result, "-----------------")
			}
		}
	} else { //Контекст до равен 0 (Выводим только строку с совпадением)
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
	filename, expression := parseFilenameExpr()

	flag.BoolVar(&fn.ignoreReg, "i", false, "Ignore register")
	flag.IntVar(&fn.afterMatch, "A", 0, "Print N strings after match")
	if fn.afterMatch < 0 {
		log.Println("-A N должно получить положительное число, так что я исправлю на ноль, типо ничего не было ;)")
	}
	flag.Parse()

	fmt.Printf("Имя файла: %s\nВыражение для поиска: %s\n", filename, expression)

	lines, err := readLines(filename)
	if err != nil {
		log.Fatal(err)
	}

	findedLines, err := fn.findStrings(lines, expression)
	if err != nil {
		log.Fatal(err)
	}
	printLines(findedLines)
}

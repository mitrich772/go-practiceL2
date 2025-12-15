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

type finder struct {
	afterMatch  int
	beforeMatch int
	ignoreReg   bool
	invert      bool
	addStrNum   bool
	onlyDigit   bool
	fixed       bool
}

type line struct {
	num  int
	text string
}

func readLines(filename string) ([]line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []line
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) != "" {
			lines = append(lines, line{
				num:  lineNum,
				text: text,
			})
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseFilenameExpr() (filename string, expr string) {
	argsLen := len(os.Args)
	expr = os.Args[argsLen-2]
	filename = os.Args[argsLen-1]
	os.Args = os.Args[:argsLen-2]
	return
}

func (fn finder) printLines(lines []line) {
	for _, l := range lines {
		if fn.addStrNum {
			fmt.Printf("%d ", l.num)
		}
		fmt.Println(l.text)
	}
}

func (fn finder) findOnlyDigit(lines []line, re *regexp.Regexp) line {
	count := 0
	for _, l := range lines {
		if fn.matchWithInvert(l, re) {
			count++
		}
	}
	return line{text: strconv.Itoa(count)}
}

func (fn finder) findWithContext(lines []line, re *regexp.Regexp) []line {
	var result []line

	for i, l := range lines {
		if !fn.matchWithInvert(l, re) {
			continue
		}

		start := i - fn.beforeMatch
		if start < 0 {
			start = 0
		}

		end := i + fn.afterMatch + 1
		if end > len(lines) {
			end = len(lines)
		}

		result = append(result, lines[start:end]...)
		result = append(result, line{text: "-----------------"})
	}

	return result
}

func (fn finder) findNoContext(lines []line, re *regexp.Regexp) []line {
	var result []line
	for _, l := range lines {
		if fn.matchWithInvert(l, re) {
			result = append(result, l)
		}
	}
	return result
}

func (fn finder) matchWithInvert(l line, re *regexp.Regexp) bool {
	return re.MatchString(l.text) != fn.invert
}

func (fn finder) findStrings(lines []line, template string) ([]line, error) {
	if fn.fixed {
		template = regexp.QuoteMeta(template)
	}
	if fn.ignoreReg {
		template = `(?i)` + template
	}

	re, err := regexp.Compile(template)
	if err != nil {
		return nil, err
	}

	switch {
	case fn.onlyDigit:
		return []line{fn.findOnlyDigit(lines, re)}, nil
	case fn.beforeMatch > 0 || fn.afterMatch > 0:
		return fn.findWithContext(lines, re), nil
	default:
		return fn.findNoContext(lines, re), nil
	}
}

func main() {
	var fn finder
	var around int

	filename, expr := parseFilenameExpr()

	flag.BoolVar(&fn.ignoreReg, "i", false, "ignore case")
	flag.BoolVar(&fn.invert, "v", false, "invert match")
	flag.BoolVar(&fn.addStrNum, "n", false, "print line number")
	flag.BoolVar(&fn.onlyDigit, "c", false, "print only match count")
	flag.BoolVar(&fn.fixed, "F", false, "fixed string search")
	flag.IntVar(&fn.afterMatch, "A", 0, "print N lines after match")
	flag.IntVar(&fn.beforeMatch, "B", 0, "print N lines before match")
	flag.IntVar(&around, "C", 0, "print N lines around match")
	flag.Parse()

	if around > 0 {
		fn.afterMatch = around
		fn.beforeMatch = around
	}

	lines, err := readLines(filename)
	if err != nil {
		log.Fatal(err)
	}

	result, err := fn.findStrings(lines, expr)
	if err != nil {
		log.Fatal(err)
	}

	fn.printLines(result)
}

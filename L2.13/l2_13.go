package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type cutter struct {
	delimiter string
	fields    string
	separated bool
}

func (ct cutter) cutStrings(lines []string) ([]string, error) {
	var res []string
	for _, line := range lines {
		parts := strings.Split(line, ct.delimiter)
		// строка без разделителя
		if len(parts) == 1 {
			if !ct.separated {
				res = append(res, line)
			}
			continue
		}
		cols := ct.parseFieldsToIntSlice(len(parts))
		parsed := ct.getNeeded(parts, cols)
		// если ничего не выбрано — пропускаем
		if parsed == "" {
			continue
		}
		res = append(res, parsed)
	}
	return res, nil
}

func (ct cutter) parseFieldsToIntSlice(max int) []int {
	var res []int

	if ct.fields == "" {
		for i := 1; i <= max; i++ {
			res = append(res, i)
		}
		return res
	}

	parts := strings.Split(ct.fields, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			if len(bounds) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start > end {
				continue
			}
			for i := start; i <= end; i++ {
				res = append(res, i)
			}
		} else {
			if n, err := strconv.Atoi(part); err == nil {
				res = append(res, n)
			}
		}
	}
	return res
}

// getNeeded возвращает строку, состоящую только из указанных колонок.
// Поля, выходящие за пределы, игнорируются (не считаются ошибкой).
func (ct cutter) getNeeded(parts []string, columnsToSave []int) string {
	var out []string
	for _, col := range columnsToSave {
		if col <= 0 || col > len(parts) {
			// игнорируем выход за пределы
			continue
		}
		out = append(out, parts[col-1])
	}
	return strings.Join(out, ct.delimiter)
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
		// сохраняем даже пустые строки — поведение можно изменить при необходимости
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}

func main() {
	var ct cutter

	flag.StringVar(&ct.delimiter, "d", "\t", "delimiter (default is tab)")
	flag.StringVar(&ct.fields, "f", "", "fields to select, e.g. 1,3-5")
	flag.BoolVar(&ct.separated, "s", false, "only print lines that contain the delimiter")
	flag.Parse()

	args := flag.Args()
	var (
		filename string
		lines    []string
		err      error
	)

	if len(args) > 0 {
		filename = args[0]
		fmt.Printf("Имя файла: %s\n", filename)
		lines, err = readLines(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// читаем из STDIN
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	cutedLines, err := ct.cutStrings(lines)
	if err != nil {
		log.Fatal(err)
	}
	printLines(cutedLines)
}

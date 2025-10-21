package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cutter struct {
	delimiter string
	fields    string
	separated bool
}

func (ct Cutter) CutStrings(lines []string) ([]string, error) {
	var res []string
	for _, line := range lines {
		cutedLine := strings.Split(line, ct.delimiter)
		if len(cutedLine) == 1 { // без разделителя строка
			if !ct.separated { // если выводим все строки
				res = append(res, cutedLine...)
			}
		} else { // с разделителем
			parsedColumnLine, err := ct.getNeeded(cutedLine, ct.parseFieldsToIntSlice(len(cutedLine)))
			if err != nil {
				return nil, err
			}
			res = append(res, parsedColumnLine)
		}
	}
	return res, nil
}

func (ct Cutter) parseFieldsToIntSlice(max int) []int {
	var res []int

	if ct.fields == "" {
		for i := 1; i <= max; i++ {
			res = append(res, i)
		}
		return res
	}

	parts := strings.Split(ct.fields, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
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

// Возвращает строку состоящую тольок из нужных колонок
func (ct Cutter) getNeeded(cutedLine []string, columnsToSave []int) (string, error) {
	var res []string
	for _, columnNumber := range columnsToSave {
		if columnNumber <= 0 || columnNumber > len(cutedLine) {
			return "", errors.New(fmt.Sprintf("Колонка %d не может быть в строке разбитой на %d колонок", columnNumber, len(cutedLine)))
		}
		res = append(res, cutedLine[columnNumber-1])
	}
	return strings.Join(res, ct.delimiter), nil
}

// Читает все непустые строки с файла
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
func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}

func main() {
	var ct Cutter

	filename := parseFilename()
	fmt.Printf("Имя файла: %s\n", filename)

	flag.StringVar(&ct.delimiter, "d", "\t", "set delimiter to strings")
	flag.StringVar(&ct.fields, "f", "", "set delimiter to strings")
	flag.BoolVar(&ct.separated, "s", false, "print only string that contain delimiter")
	flag.Parse()

	lines, err := ReadLines(filename)
	if err != nil {
		log.Fatal(err)
	}
	cutedLines, err := ct.CutStrings(lines)
	if err != nil {
		log.Println(err)
	}
	printLines(cutedLines)

}

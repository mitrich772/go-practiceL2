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

// Months содержит соответствие месяца и его номера
var Months = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}

// Sorter сортирует строки по заданным флагам
type Sorter struct {
	Column        int
	Numeric       bool
	Reverse       bool
	RemoveTBlanks bool
	Unique        bool
	CheckSort     bool
	HumanReadable bool
	MonthCheck    bool
	SortType      bool
	Err           error
}

// Хранит строку и ее числовое значение, если возможно
type sortableLine struct {
	line   string
	keyInt int
	keyStr string
}

// Преобразовывает строки в sortableLine, находя ключи
func (s Sorter) buildSortableLines(lines []string) ([]sortableLine, error) {
	res := make([]sortableLine, len(lines))

	for i, line := range lines {
		var keyInt int
		var keyStr string

		// Берем нужную колонку
		colText, err := getColumn(line, s.Column)
		if err != nil {
			return nil, err
		}

		// Если нужно, убираем хвостовые пробелы сразу
		if s.RemoveTBlanks {
			colText = strings.TrimRight(colText, " \t")
		}

		if s.Numeric {
			keyInt, err = strconv.Atoi(colText)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования '%s': %w", colText, err)
			}
		} else if s.HumanReadable {
			keyInt, err = toHumanFormat(colText)
			if err != nil {
				return nil, fmt.Errorf("ошибка преобразования '%s': %w", colText, err)
			}
		} else if s.MonthCheck {
			if !isMonth(colText) {
				return nil, fmt.Errorf("строка не месяц '%s'", colText)
			}
			keyInt = Months[colText]
		} else {
			keyStr = colText // для обычных строк берём именно колонку
		}

		res[i] = sortableLine{
			line:   line,
			keyInt: keyInt,
			keyStr: keyStr,
		}
	}

	return res, nil
}

// Sort делает сортировку выбранным методом
func (s *Sorter) Sort(lines []string) error {
	if s.SortType {
		return s.SortA(lines)
	}
	return s.SortB(lines)
}

// SortA делает сортировку полученных строк, с помощью buildSortableLines
func (s *Sorter) SortA(lines []string) error {
	s.Err = nil

	SLines, err := s.buildSortableLines(lines)
	if err != nil {
		s.Err = err
		return err
	}

	slices.SortFunc(SLines, s.compareLinesA)

	for i := range lines {
		lines[i] = SLines[i].line
	}

	return s.Err
}

// SortB делает сортировку полученных строк
func (s *Sorter) SortB(lines []string) error {
	s.Err = nil
	slices.SortFunc(lines, s.compareLinesB)
	return s.Err
}

// Функция сравнения для sortableLine
func (s *Sorter) compareLinesA(a, b sortableLine) int {
	if s.Numeric || s.HumanReadable || s.MonthCheck {
		return compareInts(a.keyInt, b.keyInt, s.Reverse)
	}

	aStr := a.keyStr
	bStr := b.keyStr

	if s.RemoveTBlanks {
		aStr = strings.TrimRight(aStr, " \t")
		bStr = strings.TrimRight(bStr, " \t")
	}

	return compareStrings(aStr, bStr, s.Reverse)
}

// Функция сравнения двух строк по заданным параметрам
func (s *Sorter) compareLinesB(a, b string) int {
	if s.Err != nil {
		return 0
	}

	// Убираем хвостовые пробелы по флагу
	if s.RemoveTBlanks {
		a = strings.TrimRight(a, " \t")
		b = strings.TrimRight(b, " \t")
	}

	// Берем нужную колонку
	var err error
	va, err := getColumn(a, s.Column)
	if err != nil {
		s.Err = err
		return 0
	}
	vb, err := getColumn(b, s.Column)
	if err != nil {
		s.Err = err
		return 0
	}

	// Сортируем как месяца
	if s.MonthCheck {
		if !isMonth(va) {
			s.Err = fmt.Errorf("строка не месяц '%s'", va)
			return 0
		}
		if !isMonth(vb) {
			s.Err = fmt.Errorf("строка не месяц '%s'", vb)
			return 0
		}

		return compareMonth(va, vb, s.Reverse)
	}
	// Сортируем по числовому значению с учётом суффиксов
	if s.HumanReadable {
		var err error
		var aInt, bInt int

		aInt, err = toHumanFormat(va)
		if err != nil {
			s.Err = fmt.Errorf("ошибка преобразования '%s': %w", va, err)
			return 0
		}

		bInt, err = toHumanFormat(vb)
		if err != nil {
			s.Err = fmt.Errorf("ошибка преобразования '%s': %w", vb, err)
			return 0
		}

		return compareInts(aInt, bInt, s.Reverse)
	}
	// Сортируем по числовому значению
	if s.Numeric {
		aInt, err := strconv.Atoi(va)
		if err != nil {
			s.Err = fmt.Errorf("ошибка преобразования '%s': %w", va, err)
			return 0
		}

		bInt, err := strconv.Atoi(vb)
		if err != nil {
			s.Err = fmt.Errorf("ошибка преобразования '%s': %w", vb, err)
			return 0
		}

		return compareInts(aInt, bInt, s.Reverse)
	}

	return compareStrings(va, vb, s.Reverse)
}

// Проверяет, отсортирован ли массив строк в соответствии с compareLines
func (s *Sorter) isSorted(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		if s.compareLinesB(lines[i-1], lines[i]) > 0 {
			return false
		}
	}
	return true
}

// Берет колонку из строки по номеру разделитель - табуляция
func getColumn(l string, column int) (string, error) {
	afterLastTabPos := 0
	curColumn := 1
	for i, v := range l {
		if v == '\t' {
			if curColumn == column {
				return l[afterLastTabPos:i], nil
			}
			curColumn++
			afterLastTabPos = i + 1
		}
	}

	if curColumn == column { //Последний столбец
		return l[afterLastTabPos:], nil
	}
	return "", fmt.Errorf("строка имеет меньше столбцов, чем k=%d: '%s'", column, l)
}

// Сравнивает строки с флагами
func compareStrings(va string, vb string, reverse bool) int {
	if reverse { //Если реверс
		return strings.Compare(vb, va)
	}
	return strings.Compare(va, vb)
}

// Сравнение чисел с реверсом
func compareInts(a, b int, reverse bool) int {
	if a == b {
		return 0
	}
	if reverse {
		if a < b {
			return 1
		}
		return -1
	}
	if a < b {
		return -1
	}
	return 1
}

// Проверяет является ли строка месяцом
func isMonth(a string) (ok bool) {
	_, ok = Months[a]
	return
}

// Сравнение строк типа месяц
func compareMonth(a, b string, reverse bool) int {
	return compareInts(Months[a], Months[b], reverse)
}

// Для первода строк вида 2к в строку 2048 уже сразу как число!
func toHumanFormat(s string) (int, error) {
	sizeSuff := map[rune]int{
		'K': 1024,
		'M': 1024 * 1024,
		'G': 1024 * 1024 * 1024,
		'k': 1024,
		'm': 1024 * 1024,
		'g': 1024 * 1024 * 1024,
	}

	multiplier := 1
	var numPart string
	var suff rune

	for i, r := range s {
		if _, ok := sizeSuff[r]; ok {
			numPart = s[:i]
			suff = r
			break
		}
	}

	if suff != 0 {
		multiplier = sizeSuff[suff]
	} else {
		numPart = s
	}

	num, err := strconv.Atoi(numPart)
	if err != nil {
		return 0, err
	}

	res := num * multiplier
	return res, nil
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

// Урезает хвостовые пробелы если надо
// func cleanTrailingBlanks(lines []string) []string {
// 	cleaned := make([]string, 0, len(lines))
// 	for _, l := range lines {
// 		s := strings.TrimRight(l, " \t") // убираем и хвостовые, и ведущие пробелы
// 		cleaned = append(cleaned, s)
// 	}
// 	return cleaned
// }

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

// Выводит массив строк
func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}

// Пишет строки в йказанный файл
func writeLinesToFile(lines []string, filename string) error {
	file, err := os.Create(filename) // Создаёт файл, перезаписывая
	if err != nil {
		return err
	}
	defer file.Close()

	for _, v := range lines {
		_, err := file.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// Разбивает слипшиеся флаги
func preprocessFlags() {
	// Слайс для новых флагов
	var expanded []string
	for _, arg := range os.Args[1:] {
		// Аргумент вида -nrbu
		if strings.HasPrefix(arg, "-") && len(arg) > 2 { // Разбиваем каждый символ в отдельный флаг
			for _, r := range arg[1:] {
				expanded = append(expanded, "-"+string(r))
			}
		} else { // Обычный флаг
			expanded = append(expanded, arg)
		}
	}
	// Заменяем os.Args на expanded
	os.Args = append([]string{os.Args[0]}, expanded...)
}

func main() {
	preprocessFlags()

	var s Sorter
	var resultToFile bool

	flag.IntVar(&s.Column, "k", 1, "column for sort")
	flag.BoolVar(&s.Numeric, "n", false, "sort string as number")
	flag.BoolVar(&s.Reverse, "r", false, "reverse sort")
	flag.BoolVar(&s.Unique, "u", false, "only unique values")
	flag.BoolVar(&s.RemoveTBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&s.CheckSort, "c", false, "check sort")
	flag.BoolVar(&s.HumanReadable, "h", false, "enable human-readable sort")
	flag.BoolVar(&s.MonthCheck, "M", false, "sort month format")
	flag.BoolVar(&s.SortType, "t", true, "swich Sort func, default SortA()")
	flag.BoolVar(&resultToFile, "f", false, "write result of sort to result_ + filename")
	flag.Parse()

	args := flag.Args() // Получаем имя файла
	var filename string
	if len(args) == 0 {
		fmt.Println("не указан файл по стандарту data.txt")
		filename = "data.txt"
	} else {
		filename = args[0]
	}

	fmt.Printf(
		"Файл: %s\ncolumn=%d, numeric=%t, reverse=%t, unique=%t, removeTBlanks=%t, checkSort=%t, resultToFile=%t\n",
		filename,
		s.Column,
		s.Numeric,
		s.Reverse,
		s.Unique,
		s.RemoveTBlanks,
		s.CheckSort,
		resultToFile,
	)

	lines, err := readLines(filename)
	if err != nil {
		log.Println(err)
	}

	if s.CheckSort {
		if s.isSorted(lines) {
			fmt.Println("-c : Строки отсортированны")
			return
		}
		fmt.Println("-c : Строки не отсортированны!")
		return
	}

	if err := s.Sort(lines); err != nil { //Если сортировка выкинула ошибку
		log.Fatalf("ошибка сортировки: %v", err)
	}

	if s.Unique {
		lines = removeDuplicatesSorted(lines)
	}
	if resultToFile { // Выбор куда выводить результат
		writeLinesToFile(lines, "result_"+filename)
	} else {
		printLines(lines)
	}
	fmt.Print("Отсортировалось")
}

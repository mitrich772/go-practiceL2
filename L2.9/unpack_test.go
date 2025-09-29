package main

import (
	"testing"
)

type test struct {
	input    string
	expected string
	wantErr  bool
}

// Примеры работы функции:

//     Вход: "a4bc2d5e"
//     Выход: "aaaabccddddde"

//     Вход: "abcd"
//     Выход: "abcd" (нет цифр — ничего не меняется)

//     Вход: "45"
//     Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)

//     Вход: ""
//     Выход: "" (пустая строка -> пустая строка)

// Дополнительное задание

// Поддерживать escape-последовательности вида \:

//     Вход: "qwe\4\5"
//     Выход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)

// Вход: "qwe\45"
// Выход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)
func UnpackTest(t *testing.T) {
	tests := []test{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
	}

	for _, curTest := range tests {
		res, err := unpack(curTest.input)
		if res != curTest.expected {
			t.Errorf("unpack(%q), expected: %q, result : %q", curTest.input, curTest.expected, res)
		}
		if (err != nil) != curTest.wantErr {
			t.Errorf("unpack(%q), wantErr: %v, error: %v", curTest.input, curTest.wantErr, err)
		}
	}
}

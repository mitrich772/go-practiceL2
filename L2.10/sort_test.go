package main

import (
	"fmt"
	"testing"
)

func TestSortByColumn(t *testing.T) {
	input := []string{
		"b\t2",
		"a\t1",
		"c\t3",
	}
	want := []string{
		"a\t1",
		"b\t2",
		"c\t3",
	}

	s := Sorter{Column: 2} // сортируем по 2 колонке
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("at %d expected %q, got %q", i, want[i], input[i])
		}
	}
}

func TestNumericSort(t *testing.T) {
	input := []string{"10", "2", "1"}
	want := []string{"1", "2", "10"}

	s := Sorter{Numeric: true, Column: 1}
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("expected %v, got %v", want, input)
		}
	}
}

func TestReverseSort(t *testing.T) {
	input := []string{"a", "b", "c"}
	want := []string{"c", "b", "a"}

	s := Sorter{Reverse: true, Column: 1}
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("expected %v, got %v", want, input)
		}
	}
}

func TestUnique(t *testing.T) {
	input := []string{"a", "a", "b", "b", "c"}
	want := []string{"a", "b", "c"}

	got := removeDuplicatesSorted(input)
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expected %v, got %v", want, got)
		}
	}
}

func TestMonthSort(t *testing.T) {
	input := []string{"Mar", "Jan", "Feb"}
	want := []string{"Jan", "Feb", "Mar"}

	s := Sorter{MonthCheck: true, Column: 1}
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("expected %v, got %v", want, input)
		}
	}
}

func TestIgnoreTrailingBlanks(t *testing.T) {
	input := []string{"a   ", "a  ", "a "}
	want := []string{"a   ", "a  ", "a "}
	s := Sorter{RemoveTBlanks: true, Column: 1}
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("expected %v, got %v", want, input)
		}
	}
}

func TestHumanReadable(t *testing.T) {
	input := []string{"1K", "512", "2K"}
	want := []string{"512", "1K", "2K"}

	s := Sorter{HumanReadable: true, Column: 1}
	if s.Sort(input); s.Err != nil {
		t.Fatalf("unexpected error: %v", s.Err)
	}

	for i := range want {
		if input[i] != want[i] {
			t.Errorf("expected %v, got %v", want, input)
		}
	}
}

func TestCheckSort(t *testing.T) {
	sorted := []string{"a", "b", "c"}
	unsorted := []string{"b", "a", "c"}

	s := Sorter{Column: 1}
	if !s.isSorted(sorted) {
		t.Error("expected sorted slice to be sorted")
	}

	if s.isSorted(unsorted) {
		t.Error("expected unsorted slice to be detected")
	}
}

func TestCheckSortByColumn(t *testing.T) {
	sorted := []string{
		"a\t1",
		"b\t2",
		"c\t3",
	}
	unsorted := []string{
		"b\t2",
		"a\t1",
		"c\t3",
	}

	// Проверяем по 2 колонке
	s := Sorter{Column: 2}

	if !s.isSorted(sorted) {
		t.Error("expected slice sorted by 2nd column to be detected as sorted")
	}

	if s.isSorted(unsorted) {
		t.Error("expected slice unsorted by 2nd column to be detected as unsorted")
	}
}

func TestSortWithInvalidColumn(t *testing.T) {
	input := []string{"a\tb", "c"}
	s := Sorter{Column: 2}
	err := s.Sort(input)
	if err == nil && s.Err == nil {
		t.Error("expected error due to missing column")
	}
}

func BenchmarkSorterNumeric(b *testing.B) {
	lines := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		lines[i] = fmt.Sprintf("%d", 1000000-i)
	}

	s := Sorter{Numeric: true}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = s.Sort(lines)
	}
}

func generateHumanReadableLines(n int) []string {
	units := []string{"B", "K", "M", "G"}
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = fmt.Sprintf("%d%s", (n-i)*10, units[i%len(units)])
	}
	return lines
}

func BenchmarkSorterHumanReadable(b *testing.B) {
	lines := generateHumanReadableLines(1000000)

	s := Sorter{HumanReadable: true}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = s.Sort(lines)
	}
}

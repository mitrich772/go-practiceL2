package main

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var testLines = []Line{
	{1, "Hello"},
	{2, "HELLO"},
	{3, "hello"},
	{4, "HeLLo"},
	{5, "f131ghello"},
	{6, "sdfshelloHElLo"},
	{7, "preHELLOpost"},
	{8, "heLLo123"},
	{9, "123HELlo456"},
	{10, "aHELloZ"},
	{11, "sayinghelloagain"},
	{12, "HELLOthere"},
	{13, "startHELLO"},
	{14, "midhelloend"},
	{15, "hellish"},
	{16, "hell"},
	{17, "the word hello is here"},
	{18, "HELLO_WORLD"},
	{19, "not this one"},
	{20, "finalHELLOtest"},
}

// helper: объединяем все строки для удобного поиска подстрок
func outputContainsAll(out []Line, subs []string) bool {
	var b strings.Builder
	for _, l := range out {
		b.WriteString(l.Text)
		b.WriteRune('\n')
	}
	joined := b.String()
	for _, s := range subs {
		if !strings.Contains(joined, s) {
			return false
		}
	}
	return true
}

func TestBasicSearch(t *testing.T) {
	fn := Finder{}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	want := []Line{
		{3, "hello"},
		{5, "f131ghello"},
		{6, "sdfshelloHElLo"},
		{11, "sayinghelloagain"},
		{14, "midhelloend"},
		{17, "the word hello is here"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("basic search mismatch.\nGot:\n%v\nWant:\n%v", got, want)
	}
}

func TestIgnoreCase(t *testing.T) {
	fn := Finder{ignoreReg: true}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{"Hello", "HELLO", "f131ghello", "HELLO_WORLD"}) {
		t.Errorf("ignore case didn't match expected lines.\nGot:\n%v", got)
	}
}

func TestInvertMatch(t *testing.T) {
	fn := Finder{invert: true}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	for _, line := range got {
		if matched, _ := regexp.MatchString("hello", line.Text); matched {
			t.Errorf("invert mode returned matching line: %s", line.Text)
		}
	}
}

func TestInvertIgnoreCase(t *testing.T) {
	fn := Finder{invert: true, ignoreReg: true}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	expected := []Line{
		{15, "hellish"},
		{16, "hell"},
		{19, "not this one"},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invert+ignore mismatch.\nGot:\n%v\nWant:\n%v", got, expected)
	}
}

func TestAfterContext(t *testing.T) {
	fn := Finder{afterMatch: 1}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{"hello", "HeLLo", "-----------------"}) {
		t.Errorf("after context (-A) missing expected parts.\nGot:\n%v", got)
	}
}

func TestBeforeContext(t *testing.T) {
	fn := Finder{beforeMatch: 1}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{"HELLO", "hello", "-----------------"}) {
		t.Errorf("before context (-B) missing expected parts.\nGot:\n%v", got)
	}
}

func TestAroundContext(t *testing.T) {
	fn := Finder{afterMatch: 1, beforeMatch: 1}
	got, err := fn.FindStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("FindStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{"HELLO", "hello", "HeLLo", "-----------------"}) {
		t.Errorf("around context (-C) missing expected parts.\nGot:\n%v", got)
	}
}

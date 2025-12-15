package main

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var testLines = []line{
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

// helper
func outputContainsAll(out []line, subs []string) bool {
	var b strings.Builder
	for _, l := range out {
		b.WriteString(l.text)
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
	fn := finder{}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	want := []line{
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
	fn := finder{ignoreReg: true}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{
		"Hello", "HELLO", "f131ghello", "HELLO_WORLD",
	}) {
		t.Errorf("ignore case didn't match expected lines.\nGot:\n%v", got)
	}
}

func TestInvertMatch(t *testing.T) {
	fn := finder{invert: true}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	for _, l := range got {
		matched, _ := regexp.MatchString("hello", l.text)
		if matched {
			t.Errorf("invert mode returned matching line: %s", l.text)
		}
	}
}

func TestInvertIgnoreCase(t *testing.T) {
	fn := finder{invert: true, ignoreReg: true}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	expected := []line{
		{15, "hellish"},
		{16, "hell"},
		{19, "not this one"},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invert+ignore mismatch.\nGot:\n%v\nWant:\n%v", got, expected)
	}
}

func TestAfterContext(t *testing.T) {
	fn := finder{afterMatch: 1}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{
		"hello", "HeLLo", "-----------------",
	}) {
		t.Errorf("after context (-A) missing expected parts.\nGot:\n%v", got)
	}
}

func TestBeforeContext(t *testing.T) {
	fn := finder{beforeMatch: 1}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{
		"HELLO", "hello", "-----------------",
	}) {
		t.Errorf("before context (-B) missing expected parts.\nGot:\n%v", got)
	}
}

func TestAroundContext(t *testing.T) {
	fn := finder{afterMatch: 1, beforeMatch: 1}
	got, err := fn.findStrings(testLines, "hello")
	if err != nil {
		t.Fatalf("findStrings failed: %v", err)
	}

	if !outputContainsAll(got, []string{
		"HELLO", "hello", "HeLLo", "-----------------",
	}) {
		t.Errorf("around context (-C) missing expected parts.\nGot:\n%v", got)
	}
}

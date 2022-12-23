package utils

import (
	"os"
	"strings"

	"golang.org/x/term"
)

type StringUtils struct{}

func (StringUtils) Reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func (StringUtils) FilterAlphaNumeric(str string) string {
	s := []byte(str)
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func (StringUtils) Flip(text string, unFlip bool) string {
	flipChars := [][]string{
		{"a", "ɐ"},
		{"b", "q"},
		{"c", "ɔ"},
		{"d", "p"},
		{"e", "ǝ"},
		{"f", "ɟ"},
		{"g", "ƃ"},
		{"h", "ɥ"},
		{"i", "ᴉ"},
		{"j", "ɾ"},
		{"k", "ʞ"},
		{"m", "ɯ"},
		{"n", "u"},
		{"r", "ɹ"},
		{"t", "ʇ"},
		{"v", "ʌ"},
		{"w", "ʍ"},
		{"y", "ʎ"},
		{"A", "∀"},
		{"C", "Ɔ"},
		{"E", "Ǝ"},
		{"F", "Ⅎ"},
		{"G", "פ"},
		{"H", "H"},
		{"I", "I"},
		{"J", "ſ"},
		{"L", "˥"},
		{"M", "W"},
		{"N", "N"},
		{"P", "Ԁ"},
		{"T", "┴"},
		{"U", "∩"},
		{"V", "Λ"},
		{"Y", "⅄"},
		{"1", "Ɩ"},
		{"2", "ᄅ"},
		{"3", "Ɛ"},
		{"4", "ㄣ"},
		{"5", "ϛ"},
		{"6", "9"},
		{"7", "ㄥ"},
		{"8", "8"},
		{"9", "6"},
		{"0", "0"},
		{".", "˙"},
		{",", "\x27"},
		{"\x27", ","},
		{"\x22", ",,"},
		{"`", ","},
		{"?", "¿"},
		{"!", "¡"},
		{"[", "]"},
		{"]", "["},
		{"(", ")"},
		{")", "("},
		{"{", "}"},
		{"}", "{"},
		{"<", ">"},
		{">", "<"},
		{"&", "⅋"},
		{"_", "‾"},
		{"∴", "∵"},
		{"⁅", "⁆"},
	}
	if unFlip {
		for _, element := range flipChars {
			text = strings.ReplaceAll(text, element[1], element[0])
		}
	} else {
		for _, element := range flipChars {
			text = strings.ReplaceAll(text, element[0], element[1])
		}
	}
	return String.Reverse(text)
}

func (StringUtils) RemoveWords(text string, amount int) string {
	words := strings.Split(text, " ")
	if amount > len(words) {
		return ""
	}
	return strings.Join(words[amount:], " ")
}

func (StringUtils) Capitalize(text string) string {
	return strings.ToUpper(text[:1]) + text[1:]
}

func (StringUtils) CenterWidth(text string, width int) string {
	spaceCount := (width - len(text)) / 2
	return strings.Repeat(" ", spaceCount) + text
}

func (StringUtils) Center(text string) string {
	fd := int(os.Stdout.Fd())
	w, _, err := term.GetSize(fd)
	if err != nil {
		return text
	}
	return String.CenterWidth(text, w)
}


func (StringUtils) CenterSlices(rows []string) string {
	fd := int(os.Stdout.Fd())
	w, _, err := term.GetSize(fd)
	if err != nil {
		return strings.Join(rows, "\n")
	}

	// Get the max length of the rows
	max := 0
	for _, row := range rows {
		if len(row) > max {
			max = len(row)
		}
	}

	// Center the rows
	spaceCount := (w - max) / 2
	spaces := strings.Repeat(" ", spaceCount)

	for i, row := range rows {
		rows[i] = spaces + row
	}

	return strings.Join(rows, "\n")
}

var String = &StringUtils{}

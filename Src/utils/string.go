package utils

import (
	"math"
	"os"
	"regexp"
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

func (StringUtils) Mock(text string) string {
	mock := ""
	for i, char := range text {
		if i%2 == 0 {
			mock += strings.ToUpper(string(char))
		} else {
			mock += strings.ToLower(string(char))
		}
	}
	return mock
}

func (StringUtils) CenterWidth(text string, width int) string {
	spaceCount := (width - len(text)) / 2
	spaceCount = int(math.Max(float64(spaceCount), 0))

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
	spaceCount = int(math.Max(float64(spaceCount), 0))

	spaces := strings.Repeat(" ", spaceCount)

	for i, row := range rows {
		rows[i] = spaces + row
	}

	return strings.Join(rows, "\n")
}

func (StringUtils) RegexReplace(regex, replaceWith string, s *string) {
	r := regexp.MustCompile(regex)

	*s = r.ReplaceAllString(*s, replaceWith)
}


func (StringUtils) RemoveDuplicateChars(word string) string {
	var last rune

	return strings.Map(func(r rune) rune {
		if r != last {
			last = r

			return r
		}

		// Negative value discards the letter
		return -1
	}, word)
}

func (StringUtils) ParseErhmahgerd(word string) string {
	// Word is too short to translate
	if len(word) < 2 {
		return word
	}

	// Common words that already have a direct translation
	switch word {
	case "AWESOME":
		return "ERSUM"
	case "BANANA":
		return "BERNERNER"
	case "BAYOU":
		return "BERU"
	case "FAVORITE", "FAVOURITE":
		return "FRAVRIT"
	case "GOOSEBUMPS":
		return "GERSBERMS"
	case "LONG":
		return "LERNG"
	case "MY":
		return "MAH"
	case "THE":
		return "DA"
	case "THEY":
		return "DEY"
	case "WE'RE":
		return "WER"
	case "YOU":
		return "U"
	case "YOU'RE":
		return "YER"
	}

	original := word

	// Remove vowels that occur at the end of the word
	// Only for words whose length is greater than 2 to prevent single character words
	if len(original) > 2 {
		String.RegexReplace(`[AEIOU]$`, "", &word)
	}

	// Reduce duplicate letters
	word = String.RemoveDuplicateChars(word)

	// Reduce consecutive vowels (and Y) to just one
	String.RegexReplace(`[AEIOUY]{2,}`, "E", &word)

	// Retain a single 'Y'
	String.RegexReplace(`Y{2,}`, "Y", &word)

	// DOWN -> DERN
	String.RegexReplace(`OW`, "ER", &word)

	// PANCAKES -> PERNKERKS
	String.RegexReplace(`AKES`, "ERKS", &word)

	// Replace vowels (and Y) with ER
	String.RegexReplace(`[AEIOUY]`, "ER", &word)

	// Other conversions that are direct
	String.RegexReplace(`ERH`, "ER", &word)
	String.RegexReplace(`MER`, "MAH", &word)
	String.RegexReplace(`ERNG`, "IN", &word)
	String.RegexReplace(`ERPERD`, "ERPED", &word)
	String.RegexReplace(`MAHM`, "MERM", &word)

	// If the word begins with Y, retain it
	if original[0] == 89 {
		word = "Y" + word
	}

	// Reduce any duplicate letters (again)
	word = String.RemoveDuplicateChars(word)

	r := regexp.MustCompile(`LOW$`)
	l := regexp.MustCompile(`LER$`)

	if r.MatchString(original) && l.MatchString(word) {
		String.RegexReplace(`LER`, "LO", &word)
	}

	return word
}

func (StringUtils) Ermahgerd(text string) string {
	var translatedWords []string
	prefix := regexp.MustCompile(`^\W+`)
	suffix := regexp.MustCompile(`\W+$`)

	text = strings.ToUpper(text)

	words := strings.Split(text, " ")

	for _, word := range words {
		wordCopy := word

		String.RegexReplace(`^\W+`, "", &word)
		String.RegexReplace(`\W+$`, "", &word)

		if len(word) > 0 {
			// Ensure that we do not lose out on any non alphabet character
			// (exclamation, question marks etc)
			beginString := prefix.FindAllString(wordCopy, 1)
			endString := suffix.FindAllString(wordCopy, 1)

			word = String.ParseErhmahgerd(word)

			if beginString != nil {
				word = beginString[0] + word
			}

			if endString != nil {
				word = word + endString[0]
			}
		} else {
			word = String.ParseErhmahgerd(word)
		}

		translatedWords = append(translatedWords, word)
	}

	return strings.Join(translatedWords, " ")
}

var String = &StringUtils{}

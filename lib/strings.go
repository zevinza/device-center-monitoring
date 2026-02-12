package lib

import (
	"regexp"
	"strings"
	"unicode"
)

// AppendStr append str with separator
func AppendStr(sep string, str ...string) string {
	return strings.Join(str, sep)
}

// StrAbbr make abbrebration of complex string
func StrAbbr(s string, length ...int) string {
	sl := strings.Split(TrimSpace(s), " ")

	r := ""
	for _, slc := range sl {
		f := slc[0]
		r += string(f)
	}

	return strings.ToUpper(r)
}

// TrimSpace remove whitespace and double space
func TrimSpace(str string) string {
	split := strings.Split(str, " ")

	var res []string
	for _, s := range split {
		if trimmed := strings.TrimSpace(s); IsValid(trimmed) {
			res = append(res, trimmed)
		}
	}
	return strings.Join(res, " ")
}

func ToSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

func RemoveSpecialChars(str string) string {
	var specialCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return specialCharsRegex.ReplaceAllString(str, "")
}

// ToTitle converts a string to title case
func ToTitle(s string) string {
	char := []rune(s)
	for i, c := range char {
		if i == 0 || unicode.IsSpace(char[i-1]) {
			char[i] = unicode.ToUpper(c)
		} else {
			char[i] = unicode.ToLower(c)
		}
	}
	return string(char)
}

// SeparateName func
func SeparateName(fullName string) (firstName, middleName, lastName string) {
	nameParts := strings.Fields(fullName)

	firstName = nameParts[0]
	lastName = nameParts[len(nameParts)-1]

	if len(nameParts) > 2 {
		middleName = strings.Join(nameParts[1:len(nameParts)-1], " ")
	}

	return firstName, middleName, lastName
}

// // RandomChars func
// func RandomChars(length int) string {
// 	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// 	b := make([]rune, length)
// 	for i := range b {
// 		b[i] = letters[mathRand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

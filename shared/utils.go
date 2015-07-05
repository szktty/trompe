package trompe

import (
	"os"
	"unicode"
)

func BeginsWithUpperCase(s string) bool {
	return unicode.IsUpper(rune(s[0]))
}

func BeginsWithLowerCase(s string) bool {
	return unicode.IsLower(rune(s[0]))
}

func IsSnakeCase(s string) bool {
	for i, c := range s {
		if i > 0 && unicode.IsUpper(rune(c)) {
			return false
		}
	}
	return true
}

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}

func StringIndex(list []string, s string) (int, bool) {
	for i, e := range list {
		if e == s {
			return i, true
		}
	}
	return -1, false
}

func ContainsString(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

func AppendStringIfAbsent(list []string, s string) []string {
	if !ContainsString(list, s) {
		list = append(list, s)
	}
	return list
}

func RevString(list []string) []string {
	rev := make([]string, len(list))
	for i, s := range list {
		rev[len(list)-i-1] = s
	}
	return rev
}

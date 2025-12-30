package cclUtils

import (
	"strings"

	"github.com/ALiwoto/ssg/ssg"
)

func SnakeToTitle(s string) string {
	bd := strings.Builder{}

	for _, split := range strings.Split(s, "_") {
		bd.WriteString(ssg.Title(split))
	}

	return bd.String()
}

func ToCamelCase(s string) string {
	title := SnakeToTitle(s)

	return strings.ToLower(title[:1]) + title[1:]
}

func ToPascalCase(str string) string {
	title := SnakeToTitle(str)

	return strings.ToUpper(title[:1]) + title[1:]
}

func ToSnakeCase(camel string) (snake string) {
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
		// disregard first letter
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

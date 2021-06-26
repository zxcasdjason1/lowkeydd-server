package twitch

import (
	"regexp"
	"strings"
)

func RemoveQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}
func AddComma(str string) string {
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}

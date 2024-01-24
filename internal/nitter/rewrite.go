package nitter

import (
	"regexp"
)

var re = regexp.MustCompile(`\b(x.com|twitter.com)\b`)

// Replace twitter.com and x.com with nitter.net
func Rewrite(url string) string {
	return re.ReplaceAllString(url, "nitter.net")
}

package nitter

import (
	"regexp"
)

// Replace twitter.com and x.com with nitter.net
func Rewrite(url string) string {
	re := regexp.MustCompile("\\b(x.com|twitter.com)\\b")
	return re.ReplaceAllString(url, "nitter.net")
}

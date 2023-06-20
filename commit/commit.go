package commit

import "fmt"

var Rev = "main"

func URL() string {
	return fmt.Sprintf("https://github.com/rcy/annie/commit/%s", Rev)
}

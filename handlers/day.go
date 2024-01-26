package handlers

import (
	"goirc/bot"
	"goirc/shell"
)

func NationalDay(params bot.HandlerParams) (string, error) {
	return shell.Command(`curl -s https://nationaltoday.com/ | pup 'meta[name=description] json{}' | jq -r .[].content | cut -f1 -d.`)
}

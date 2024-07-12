package day

import (
	"goirc/bot"
)

var url = "https://www.daysoftheyear.com/today/"

var dayCache = NewCache(`curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Day$'`)
var weekCache = NewCache(`curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Week$'`)
var monthCache = NewCache(`curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Month$'`)

func NationalDay(params bot.HandlerParams) error {
	str, err := dayCache.Pop()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", str)

	return nil
}

func NationalWeek(params bot.HandlerParams) error {
	str, err := weekCache.Pop()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", str)

	return nil
}

func NationalMonth(params bot.HandlerParams) error {
	str, err := monthCache.Pop()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", str)

	return nil
}

func NationalRefs(params bot.HandlerParams) error {
	params.Privmsgf(params.Target, "%s", url)

	return nil
}

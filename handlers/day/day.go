package day

import (
	"context"
	"fmt"
	"goirc/bot"
	"goirc/image"
	"goirc/shell"
	"strings"
)

var url = "https://www.daysoftheyear.com/today/"

var dayCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E '\bDay\b'`
var weekCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E '\bWeek\b'`
var monthCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E '\bMonth\b'`

var dayCache = NewCache(dayCmd)
var weekCache = NewCache(weekCmd)
var monthCache = NewCache(monthCmd)

func NationalDay(ctx context.Context, params bot.HandlerParams) error {
	str, err := dayCache.Pop()
	if err != nil {
		return err
	}

	str = strings.ReplaceAll(str, "&amp;", "&")

	if str == "EOF" {
		img, err := dayImage(ctx, dayCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "Today's image: %s", img.URL())
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalWeek(ctx context.Context, params bot.HandlerParams) error {
	str, err := weekCache.Pop()
	if err != nil {
		return err
	}

	str = strings.ReplaceAll(str, "&amp;", "&")

	if str == "EOF" {
		img, err := dayImage(ctx, weekCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "This week's image: %s", img.URL())
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalMonth(ctx context.Context, params bot.HandlerParams) error {
	str, err := monthCache.Pop()
	if err != nil {
		return err
	}

	str = strings.ReplaceAll(str, "&amp;", "&")

	if str == "EOF" {
		img, err := dayImage(ctx, monthCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "This month's image: %s", img.URL())
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalRefs(ctx context.Context, params bot.HandlerParams) error {
	params.Privmsgf(params.Target, "%s", url)

	return nil
}

func stripPhrases(days []string) []string {
	removes := []string{
		"and",
		"day",
		"for",
		"international",
		"month",
		"national",
		"the",
		"week",
		"weekend",
		"world",
		"year",
	}
	result := make([]string, len(days))

	var kept []string

	for d, day := range days {
		day = strings.ToLower(day)

		kept = []string{}

		for _, word := range strings.Fields(day) {
			keep := true

			for _, remove := range removes {
				if word == remove {
					keep = false
					break
				}
			}

			if keep {
				kept = append(kept, word)
			}
		}
		result[d] = strings.Join(kept, " ")
	}
	return result
}

func dayImage(ctx context.Context, cmd string) (*image.GeneratedImage, error) {
	r, err := shell.Command(cmd)
	if err != nil {
		return nil, err
	}

	r = strings.ReplaceAll(r, "&amp;", "&")

	days := strings.Split(strings.TrimSpace(r), "\n")
	days = stripPhrases(days)
	prompt := "create a single scene with representations of " + strings.Join(days, ", ")
	gi, err := image.GenerateDALLE(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("prompt: %s: %w", prompt, err)
	}

	return gi, nil
}

func Dayi(ctx context.Context, params bot.HandlerParams) error {
	img, err := dayImage(ctx, dayCmd)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "Today's image: %s", img.URL())
	return nil
}

func Weeki(ctx context.Context, params bot.HandlerParams) error {
	img, err := dayImage(ctx, weekCmd)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "This week's image: %s", img.URL())
	return nil
}

func Monthi(ctx context.Context, params bot.HandlerParams) error {
	img, err := dayImage(ctx, monthCmd)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "This month's image: %s", img.URL())
	return nil
}

func Image(ctx context.Context, params bot.HandlerParams) error {
	prompt := params.Matches[1]
	gi, err := image.GenerateDALLE(ctx, prompt)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", gi.URL())

	return nil
}

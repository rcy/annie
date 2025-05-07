package tip

import (
	"goirc/bot"
)

func Handle(params bot.HandlerParams) error {
	params.Privmsgf(params.Target, "%s: https://rcy.sh/tip", params.Nick)

	return nil
}

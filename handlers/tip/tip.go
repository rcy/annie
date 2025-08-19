package tip

import "goirc/internal/responder"

func Handle(params responder.Responder) error {
	params.Privmsgf(params.Target(), "%s: https://rcy.sh/tip", params.Nick)

	return nil
}

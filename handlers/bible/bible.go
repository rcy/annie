package bible

import (
	"goirc/bot"
	"goirc/internal/bible"
)

func Handle(params bot.HandlerParams) error {
	ref := params.Matches[1]

	b := bible.New()

	err := b.SetActiveTranslation("KJV")
	if err != nil {
		return err
	}

	text, err := b.Lookup(ref)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", text)

	return nil
}

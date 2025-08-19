package bible

import (
	"goirc/internal/bible"
	"goirc/internal/responder"
)

func Handle(params responder.Responder) error {
	ref := params.Match(1)

	b := bible.New()

	err := b.SetActiveTranslation("KJV")
	if err != nil {
		return err
	}

	text, err := b.Lookup(ref)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target(), "%s", text)

	return nil
}

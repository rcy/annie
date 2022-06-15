package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"fmt"
)

const channel = "#embx";
const serverssl = "irc.libera.chat:6697"

func main() {
        ircnick1 := "annie"
        irccon := irc.IRC(ircnick1, "IRCTestSSL")
        irccon.VerboseCallbackHandler = true
        irccon.Debug = true
        irccon.UseTLS = true
        irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
        irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
        irccon.AddCallback("366", func(e *irc.Event) {  })
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		irccon.Privmsg(e.Arguments[0], "echo")
	})
        err := irccon.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
        irccon.Loop()
}

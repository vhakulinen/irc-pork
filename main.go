package main

import (
	"github.com/vhakulinen/irc-pork/commands"
	"github.com/vhakulinen/irc-pork/ui"
)

var defaultNick = "Girc"

func main() {
	ui.Init()
	defer ui.Close()
	// Loop to read data from UI
	go func() {
		for {
			select {
			case msg, ok := <-ui.Input:
				if !ok {
					ui.Writer.Write([]byte("Failed to read data from ui.Input"))
					break
				}
				commands.Handle(msg)
				break
			}
		}
	}()

	ui.Loop()
}

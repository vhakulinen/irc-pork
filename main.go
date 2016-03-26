package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sorcix/irc"
	"github.com/vhakulinen/girc/ui"
)

func main() {
	ui.Init()
	defer ui.Close()

	conn, err := irc.Dial("irc.unix.chat:6667")
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
		os.Exit(2)
	}

	_, err = conn.Encoder.Write([]byte("USER k k k k"))
	_, err = conn.Encoder.Write([]byte("NICK k"))
	if err != nil {
		fmt.Printf("Failed to write: %v", err)
		os.Exit(2)
	}

	// Loop to read data from IRC
	go func() {
		for {
			msg, err := conn.Decode()
			if err != nil {
				ui.Writer.Write([]byte(fmt.Sprintf("Failed to read data from IRC: %v", err)))
				continue
			}
			if msg.Command == irc.PING {
				conn.Encoder.Write([]byte(fmt.Sprintf("PONG :%s", msg.Trailing)))
			} else if msg.Command == irc.PRIVMSG {
				ui.Write(msg.Params[0], fmt.Sprintf("%s @ %s: %s",
					msg.User, msg.Params[0], msg.Trailing))
			} else {
				ui.Writer.Write([]byte(msg.Bytes()))
			}
		}
	}()

	// Loop to read data from UI
	go func() {
		for {
			select {
			case msg, ok := <-ui.Input:
				if !ok {
					ui.Writer.Write([]byte("Failed to read data from ui.Input"))
					break
				}
				ui.Writer.Write([]byte(msg.Message))
				args := strings.Split(msg.Message, " ")
				if len(args) > 1 && args[0] == "join" {
					conn.Encoder.Write([]byte(fmt.Sprintf("JOIN :%s", args[1])))
				} else {
					conn.Encoder.Write([]byte(fmt.Sprintf("PRIVMSG %s :%s",
						msg.Target, msg.Message)))
					ui.Write(msg.Target, fmt.Sprintf("ME: %s", msg.Message))
				}
				break
			}
		}
	}()

	ui.Loop()
}

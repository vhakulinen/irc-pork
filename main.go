package main

import (
	"fmt"
	"strings"

	"github.com/sorcix/irc"
	"github.com/vhakulinen/girc/ui"
)

var defaultNick = "Girc"

func connect(addr string) {
	conn, err := irc.Dial(addr)
	if err != nil {
		ui.Writer.Write([]byte((fmt.Sprintf("Failed to connect: %v", err))))
		return
	}

	_, err = conn.Encoder.Write([]byte(fmt.Sprintf("USER %s %s %s %s",
		defaultNick, defaultNick, defaultNick, defaultNick)))
	_, err = conn.Encoder.Write([]byte(fmt.Sprintf("NICK %s", defaultNick)))
	if err != nil {
		fmt.Printf("Failed to write: %v", err)
		conn.Close()
		return
	}

	ui.Connections.AddConnection(conn)

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
					msg.User, msg.Params[0], msg.Trailing), conn)
			} else {
				ui.Writer.Write([]byte(msg.Bytes()))
			}
		}
	}()
}

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
				ui.Writer.Write([]byte(msg.Message))
				args := strings.Split(msg.Message, " ")
				if len(args) > 1 && args[0] == "/connect" {
					connect(args[1])
				} else if len(args) > 1 && args[0] == "/join" {
					if msg.Conn == nil {
						break
					}
					msg.Conn.Encoder.Write([]byte(fmt.Sprintf("JOIN :%s", args[1])))
				} else {
					if msg.Conn == nil {
						break
					}
					msg.Conn.Encoder.Write([]byte(fmt.Sprintf("PRIVMSG %s :%s",
						msg.Target, msg.Message)))
					ui.Write(msg.Target, fmt.Sprintf("ME: %s", msg.Message), msg.Conn)
				}
				break
			}
		}
	}()

	ui.Loop()
}

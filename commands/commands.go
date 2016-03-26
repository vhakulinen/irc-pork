// Package commands handles UI commands
package commands

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sorcix/irc"
	"github.com/vhakulinen/girc/ui"
	"github.com/vhakulinen/girc/utils"
)

var commands = map[string]func(args []string, target string, conn *utils.Connection,
	output io.Writer){
	// Echo echoes the data to output window
	"echo": func(args []string, target string, conn *utils.Connection, output io.Writer) {
		ui.Writer.Write([]byte(strings.Join(args[1:], " ")))
	},
	// Connect connects to server using passed irc connection
	"connect": ircCommandConnect,
	"join":    ircCommandJoin,
	"j":       ircCommandJoin,
	"privmsg": ircCommandPrivMsg,
}

// Handle handles UI input.
func Handle(input *ui.InputData) {
	if input.Message[0] != '/' {
		input.Message = "/privmsg " + input.Message
	}
	args := strings.Split(input.Message[1:], " ")
	for name, parse := range commands {
		if name == args[0] {
			parse(args, input.Target, input.Conn, input.Writer)
			break
		}
	}
}

func ircCommandJoin(args []string, target string, conn *utils.Connection, output io.Writer) {
	if len(args) == 1 {
		ui.Writer.Write([]byte("Usage"))
		ui.Writer.Write([]byte(fmt.Sprintf("\t%s <channel>", args[0])))
		return
	}
	conn.Encoder.Encode(&irc.Message{
		Command: irc.JOIN,
		Params:  args[1:],
	})
	//var conn = pool.Current
	//conn.Join(args[1])
}

func ircCommandConnect(args []string, target string, _ *utils.Connection, output io.Writer) {
	var defaultNick = "Girc"
	usage := func() {
		ui.Writer.Write([]byte("Usage:"))
		ui.Writer.Write([]byte("\tconnect <host> [<port>]"))
	}
	if len(args) == 1 {
		usage()
		return
	}
	host := args[1]
	port := 6667
	if len(args) == 3 {
		var err error
		port, err = strconv.Atoi(args[2])
		if err != nil {
			usage()
			return
		}
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	iconn, err := irc.Dial(addr)
	if err != nil {
		ui.Writer.Write([]byte((fmt.Sprintf("Failed to connect: %v", err))))
		return
	}
	conn := &utils.Connection{
		iconn, addr,
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
	// TODO: Add wait group for connections
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

func ircCommandPrivMsg(args []string, target string, conn *utils.Connection, output io.Writer) {
	err := conn.Encoder.Encode(&irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{target},
		Trailing: strings.Join(args[1:], " "),
	})
	if err != nil {
		ui.Writer.Write([]byte(fmt.Sprintf("Failed to send message: %v", err)))
	} else {
		output.Write([]byte(fmt.Sprintf("MEEE!!! @ %s: %s", target, strings.Join(args, " "))))
	}
}

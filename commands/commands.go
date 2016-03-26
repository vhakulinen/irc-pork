// commands package handles UI commands
package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sorcix/irc"
	"github.com/vhakulinen/girc/utils"
)

var commands = map[string]func(args []string, channel *utils.Channel,
	pool *utils.IrcPool, logger *log.Logger){
	// Echo echoes the data to output window
	"echo": func(args []string, channel *utils.Channel, pool *utils.IrcPool, logger *log.Logger) {
		logger.Println(strings.Join(args, " "))
	},
	// Connect connects to server using passed irc connection
	"connect": ircCommandConnect,
	"join":    ircCommandJoin,
	"j":       ircCommandJoin,
	"privmsg": ircCommandPrivMsg,
}

// Handle handles UI input.
func Handle(cmd string, channel *utils.Channel, pool *utils.IrcPool, logger *log.Logger) {
	if cmd[0] != '/' {
		cmd = "/privmsg " + cmd
	}
	args := strings.Split(cmd[1:], " ")
	for name, parse := range commands {
		if name == args[0] {
			parse(args, channel, pool, logger)
			break
		}
	}
}

func ircCommandJoin(args []string, channel *utils.Channel, pool *utils.IrcPool, logger *log.Logger) {
	if len(args) == 1 {
		logger.Println("Usage")
		logger.Println(fmt.Sprintf("\t%s <channel>", args[0]))
		return
	}
	//var conn = pool.Current
	//conn.Join(args[1])
}

func ircCommandConnect(args []string, channel *utils.Channel, pool *utils.IrcPool, logger *log.Logger) {
	usage := func() {
		logger.Println("Usage:")
		logger.Println("\tconnect <host> [<port>]")
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
	if conn, err := irc.Dial(fmt.Sprintf("%s:%d", host, port)); err == nil {
		pool.AddConnection(conn)
	} else {
		logger.Printf("Error: %v\n", err)
	}
}

func ircCommandPrivMsg(args []string, channel *utils.Channel, pool *utils.IrcPool, logger *log.Logger) {
	channel.Send(&irc.Message{
		Command: irc.PRIVMSG,
		Params:  args,
	})
}

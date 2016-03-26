package utils

import "github.com/sorcix/irc"

// Connection extends irc.Connection so that we have more data
// avaiable about it.
type Connection struct {
	*irc.Conn
	Name string
}

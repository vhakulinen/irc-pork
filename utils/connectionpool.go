package utils

import (
	"sync"

	"github.com/sorcix/irc"
)

// ConnectionPool is just a pool which contains connections (du'h).
type ConnectionPool struct {
	*sync.Mutex
	pool []Connection
}

// NewConnectionPool creates new ConnectionPool with empty pool.
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		&sync.Mutex{},
		[]*irc.Conn{},
	}
}

// AddConnection adds connection to ConnectionPool if it
// isn't there yet.
func (cp *ConnectionPool) AddConnection(conn *irc.Conn) {
	cp.Lock()
	for _, c := range cp.pool {
		if c == conn {
			break
		}
	}
	cp.pool = append(cp.pool, conn)
	cp.Unlock()
}

// GetPool returns the pool.
func (cp *ConnectionPool) GetPool() []*irc.Conn {
	return cp.pool
}

package utils

import (
	"sync"
)

// ConnectionPool is just a pool which contains connections (du'h).
type ConnectionPool struct {
	*sync.Mutex
	pool []*Connection
}

// NewConnectionPool creates new ConnectionPool with empty pool.
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		&sync.Mutex{},
		[]*Connection{},
	}
}

// AddConnection adds connection to ConnectionPool if it
// isn't there yet.
func (cp *ConnectionPool) AddConnection(conn *Connection) {
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
func (cp *ConnectionPool) GetPool() []*Connection {
	return cp.pool
}

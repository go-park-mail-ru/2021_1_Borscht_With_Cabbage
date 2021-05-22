package websocketPool

import (
	"sync"

	"github.com/gorilla/websocket"
)

func NewConnectionPool() ConnectionPool {
	return ConnectionPool{
		Connections: make(map[int]*websocket.Conn),
	}
}

type ConnectionPool struct {
	mutex       sync.RWMutex
	Connections map[int]*websocket.Conn
}

func (p *ConnectionPool) Add(id int, ws *websocket.Conn) {
	p.mutex.Lock()
	p.Connections[id] = ws
	p.mutex.Unlock()
}

func (p *ConnectionPool) Get(id int) *websocket.Conn {
	return p.Connections[id]
}

func (p *ConnectionPool) Remove(id int) {
	p.mutex.Lock()
	p.Connections[id] = nil
	p.mutex.Unlock()
}

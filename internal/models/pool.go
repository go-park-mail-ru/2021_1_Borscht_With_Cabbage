package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

func NewConnectionPool() ConnectionPool {
	return ConnectionPool{
		connections: make(map[int]*websocket.Conn),
	}
}

type ConnectionPool struct {
	mutex       sync.RWMutex
	connections map[int]*websocket.Conn
}

func (p *ConnectionPool) Add(id int, ws *websocket.Conn) {
	p.mutex.Lock()
	p.connections[id] = ws
	p.mutex.Unlock()
}

func (p *ConnectionPool) Get(id int) *websocket.Conn {
	return p.connections[id]
}

func (p *ConnectionPool) Remove(id int) {
	p.mutex.Lock()
	p.connections[id] = nil
	p.mutex.Unlock()
}

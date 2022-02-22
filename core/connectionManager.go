package core

import (
	"errors"
	"fmt"
	_interface "gotoCabbage/interface"
	"sync"
)

type ConnectionManager struct {
	conns map[uint32]_interface.IConnection
	mutex sync.RWMutex
}

func (c *ConnectionManager) Add(conn _interface.IConnection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.conns[conn.GetConnID()] = conn
	fmt.Println("connection Add success")
}

func (c *ConnectionManager) Remove(conn _interface.IConnection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.conns, conn.GetConnID())
	fmt.Println("connection Remove success")
}

func (c *ConnectionManager) Get(connID uint32) (_interface.IConnection, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if conn, ok := c.conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnectionManager) Count() int {
	return len(c.conns)
}

func (c *ConnectionManager) ClearConn() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for connID, conn := range c.conns {
		conn.Stop()
		delete(c.conns, connID)
	}
	fmt.Println("clear all Connection success")
}

func NewConnManager() _interface.IConnectionManager {
	return &ConnectionManager{
		conns: make(map[uint32]_interface.IConnection),
		mutex: sync.RWMutex{},
	}
}

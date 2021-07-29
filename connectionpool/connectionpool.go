package connectionpool

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionPool struct {
	connectionNum     int64
	connectionCreator func() (net.Conn, error)
	connections       chan net.Conn
	mu                *sync.Mutex
	isClosed          bool
}

var connectionPool *ConnectionPool

func InitConnectionPool(poolSize int64, creator func() (net.Conn, error)) {
	connectionPool = &ConnectionPool{
		connectionNum:     poolSize,
		connectionCreator: creator,
		connections:       make(chan net.Conn, poolSize),
		isClosed:          false,
		mu:                &sync.Mutex{},
	}
	for i := 0; i < int(poolSize); i++ {
		go func() {
			conn, err := connectionPool.connectionCreator()
			if err != nil {
				fmt.Println(err)
				return
			}
			connectionPool.connections <- conn
		}()
	}
}

// Get get a conn from pool, return an error if pool is closed.
func (pool *ConnectionPool) Get() (net.Conn, error) {
	var closed bool
	pool.mu.Lock()
	closed = pool.isClosed
	pool.mu.Unlock()
	if closed {
		return nil, fmt.Errorf("[%s]: %s", "Get Connection failed", "pool is already closed")
	}
	conn := <-pool.connections
	return conn, nil
}

// Put put the conn to the pool after using, return an error if pool is closed
func (pool *ConnectionPool) Put(conn net.Conn) error {
	var closed bool
	pool.mu.Lock()
	closed = pool.isClosed
	pool.mu.Unlock()
	if closed {
		return fmt.Errorf("[%s]: %s", "Put Connection failed", "pool is already closed")
	}
	pool.connections <- conn
	return nil
}

// Close it's better close it if not use.
func (pool *ConnectionPool) Close() {
	if pool.isClosed {
		return
	}
	pool.mu.Lock()
	if !pool.isClosed {
		pool.isClosed = true
		for i := 0; i < int(pool.connectionNum); i++ {
			conn := <-pool.connections
			err := conn.Close()
			if err != nil {
				//pool.errors <- err
				fmt.Println(err.Error())
				continue
			}
		}
	}
	pool.mu.Unlock()
}

func Get() (net.Conn, error) {
	if connectionPool == nil {
		return nil, fmt.Errorf("Please init pool first")
	}
	return connectionPool.Get()
}

func Put(conn net.Conn) error {
	return connectionPool.Put(conn)
}

func Close() {
	connectionPool.Close()
}

func Len() int {
	return len(connectionPool.connections)
}

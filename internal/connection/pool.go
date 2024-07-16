package connection

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const maxQueueLength = 10_000

// Manages tcp connections for publishers
type TCPPool struct {
	host     string
	port     int
	mtx      sync.Mutex
	numOpen  int
	maxOpen  int
	requests chan *request
}

func NewTCPPool(host string, port, maxOpen, maxIdle int) *TCPPool {
	pool := &TCPPool{
		host:     host,
		port:     port,
		maxOpen:  maxOpen,
		requests: make(chan *request, maxQueueLength),
	}

	go pool.handleRequests()

	return pool
}

func (p *TCPPool) Close(conn TCPConn) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	err := conn.Close()
	if err != nil {
		return err
	}
	p.numOpen--
	return nil
}

func (p *TCPPool) Get() (TCPConn, error) {
	p.mtx.Lock()

	if p.maxOpen > 0 && p.numOpen > 0 && p.numOpen >= p.maxOpen {
		req := &request{
			connections: make(chan TCPConn, 1),
			errors:      make(chan error, 1),
		}

		p.requests <- req
		p.mtx.Unlock()

		select {
		case tcpConn := <-req.connections:
			return tcpConn, nil
		case err := <-req.errors:
			return nil, err
		}
	}

	p.numOpen++
	p.mtx.Unlock()

	newConn, err := p.openNewConnection()
	if err != nil {
		p.mtx.Lock()
		p.numOpen--
		p.mtx.Unlock()
		return nil, err
	}

	return newConn, nil
}

func (p *TCPPool) openNewConnection() (TCPConn, error) {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn := NewConnection(c)
	return conn, nil
}

func (p *TCPPool) handleRequests() {
	for req := range p.requests {
		requestDone := false
		hasTimeout := false
		timeout := time.After(30 * time.Second)

		for {
			if requestDone || hasTimeout {
				break
			}
			select {
			case <-timeout:
				hasTimeout = true
				req.errors <- fmt.Errorf("connection request timed out")
			default:
				p.mtx.Lock()
				if p.maxOpen > 0 && p.numOpen < p.maxOpen {
					p.numOpen++
					p.mtx.Unlock()

					c, err := p.openNewConnection()
					if err != nil {
						p.mtx.Lock()
						p.numOpen--
						p.mtx.Unlock()
					} else {
						req.connections <- c
						requestDone = true
					}
				} else {
					p.mtx.Unlock()
				}
			}
		}

	}
}

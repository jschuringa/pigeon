package connection

import (
	"net"
)

type TCPConn interface {
	Send(data []byte) error
	Close() error
	ID() string
}

type tcpConn struct {
	id   string
	conn net.Conn
}

func NewConnection(conn net.Conn) TCPConn {
	return &tcpConn{
		conn: conn,
	}
}

func (c *tcpConn) Send(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *tcpConn) Close() error {
	return c.conn.Close()
}

func (c *tcpConn) ID() string {
	return c.id
}

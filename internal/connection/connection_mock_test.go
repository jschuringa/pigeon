package connection_test

import (
	"net"
	"time"
)

type mockNetConnection struct {
	err error
}

// LocalAddr implements net.Conn.
func (m *mockNetConnection) LocalAddr() net.Addr {
	panic("unimplemented")
}

// Read implements net.Conn.
func (m *mockNetConnection) Read(b []byte) (n int, err error) {
	return 0, m.err
}

// RemoteAddr implements net.Conn.
func (m *mockNetConnection) RemoteAddr() net.Addr {
	panic("unimplemented")
}

// SetDeadline implements net.Conn.
func (m *mockNetConnection) SetDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetReadDeadline implements net.Conn.
func (m *mockNetConnection) SetReadDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetWriteDeadline implements net.Conn.
func (m *mockNetConnection) SetWriteDeadline(t time.Time) error {
	panic("unimplemented")
}

// Write implements net.Conn.
func (m *mockNetConnection) Write(b []byte) (n int, err error) {
	return 0, m.err
}

func (m *mockNetConnection) Close() error {
	return m.err
}

var _ net.Conn = (*mockNetConnection)(nil)

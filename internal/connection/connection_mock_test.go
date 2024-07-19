package connection_test

import (
	"net"
	"time"
)

type MockNetConnection struct {
	err error
}

// LocalAddr implements net.Conn.
func (m *MockNetConnection) LocalAddr() net.Addr {
	panic("unimplemented")
}

// Read implements net.Conn.
func (m *MockNetConnection) Read(b []byte) (n int, err error) {
	return 0, m.err
}

// RemoteAddr implements net.Conn.
func (m *MockNetConnection) RemoteAddr() net.Addr {
	panic("unimplemented")
}

// SetDeadline implements net.Conn.
func (m *MockNetConnection) SetDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetReadDeadline implements net.Conn.
func (m *MockNetConnection) SetReadDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetWriteDeadline implements net.Conn.
func (m *MockNetConnection) SetWriteDeadline(t time.Time) error {
	panic("unimplemented")
}

// Write implements net.Conn.
func (m *MockNetConnection) Write(b []byte) (n int, err error) {
	return 0, m.err
}

func (m *MockNetConnection) Close() error {
	return m.err
}

var _ net.Conn = (*MockNetConnection)(nil)

package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

type Listener struct {
	host string
	port int

	lc net.ListenConfig
}

func NewListener(host string, port int) *Listener {
	l := &Listener{
		host: host,
		port: port,
	}

	l.lc = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				if err != nil {
					log.Printf("Could not set SO_REUSEADDR sockeet option: %s", err)
				}

				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				if err != nil {
					log.Printf("Could not set SO_REUSEPORT socket option: %s", err)
				}
			})
		},
	}

	return l
}

func (l *Listener) Start(ctx context.Context) (net.Listener, error) {
	conn, err := l.lc.Listen(ctx, "tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

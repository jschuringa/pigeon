package connection

type request struct {
	connections chan TCPConn
	errors      chan error
}

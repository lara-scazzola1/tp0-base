package common

import (
	"io"
	"net"
)

type Socket struct {
	conn net.Conn
}

func NewSocket(conn net.Conn) *Socket {
	return &Socket{conn: conn}
}

// Recvall reads exactly size bytes from the connected socket and
// stores them in buf. Returns an error if the reading fails.
func (skt *Socket) Recvall(size int, buf []byte) error {
	numBytesToRead := 0
	for numBytesToRead < size {
		i, err := skt.conn.Read(buf[numBytesToRead:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		numBytesToRead += i
	}
	return nil
}

// Sendall writes exactly size bytes from buf to the connected socket.
// Returns an error if the writing fails.
func (skt *Socket) Sendall(size int, buf []byte) error {
	numBytesToSend := 0
	for numBytesToSend < size {
		i, err := skt.conn.Write(buf[numBytesToSend:size])
		if err != nil {
			return err
		}
		numBytesToSend += i
	}
	return nil
}

// Close closes the connection to the socket.
func (skt *Socket) Close() error {
	return skt.conn.Close()
}

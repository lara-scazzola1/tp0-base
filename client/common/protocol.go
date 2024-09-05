package common

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	// comandos que envia
	BET_COMMAND = 0

	// comandos que recibe
	RESPONSE_BET_COMMAND = 0
)

type Protocol struct {
	socket *Socket
}

// NewProtocol creates a new Protocol instance with the provided connection.
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{socket: NewSocket(conn)}
}

// SendBet serializes and sends the provided bet data to the connected socket.
// Encodings:
// - uint8: command
// - uint8: length of the serialized bet data
// - serialized bet data
// Returns an error if the serialization or sending fails.
func (p *Protocol) SendBet(bet *Bet) error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(BET_COMMAND))

	serializeData, err := bet.Serialize()
	if err != nil {
		return err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint8(len(serializeData))); err != nil {
		return err
	}

	buffer.Write(serializeData)

	err = p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// ReceiveResponseBet reads the response command from the connected socket.
// Returns true if the response command is RESPONSE_BET_COMMAND, false otherwise.
func (p *Protocol) ReceiveResponseBet() (bool, error) {
	buf := make([]byte, 1)
	err := p.socket.Recvall(1, buf)
	if err != nil {
		return false, err
	}
	if buf[0] != RESPONSE_BET_COMMAND {
		return false, nil
	}
	return true, nil
}

// Close closes the connection of the Protocol instance.
func (p *Protocol) Close() error {
	return p.socket.Close()
}

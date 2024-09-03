package common

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	// comandos que envia
	BetCommand = 0

	// comandos que recibe
	ResponseBetCommand = 0
)

type Protocol struct {
	socket *Socket
}

func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{socket: NewSocket(conn)}
}

func (p *Protocol) SendBet(bet *Bet) error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(BetCommand))

	serializeData, err := bet.Serialize()
	if err != nil {
		return err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(serializeData))); err != nil {
		return err
	}

	buffer.Write(serializeData)

	err = p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (p *Protocol) ReceiveResponseBet() (bool, error) {
	buf := make([]byte, 1)
	err := p.socket.Recvall(1, buf)
	if err != nil {
		return false, err
	}
	if buf[0] != ResponseBetCommand {
		return false, nil
	}
	return true, nil
}

func (p *Protocol) Close() error {
	return p.socket.Close()
}

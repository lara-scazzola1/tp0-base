package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	MAX_SIZE_BATCH = 8 * 1024

	// Comandos que envia
	BATCH_COMMAND      = 1
	DISCONNECT_COMMAND = 2

	// Comandos que recibe
	RESPONSE_BATCH_COMMAND_OK    = 1
	RESPONSE_BATCH_COMMAND_ERROR = 2
)

type Protocol struct {
	socket *Socket
}

func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{socket: NewSocket(conn)}
}

func serializeCommandBet(bet *Bet) ([]byte, error) {
	buffer := new(bytes.Buffer)

	serializeData, err := bet.Serialize()
	if err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint8(len(serializeData))); err != nil {
		return nil, err
	}

	buffer.Write(serializeData)

	return buffer.Bytes(), nil
}

func (p *Protocol) SendBatch(bets []*Bet, exit chan os.Signal) error {
	bufferBetsData := new(bytes.Buffer)
	for _, bet := range bets {
		select {
		case <-exit:
			return fmt.Errorf("exit signal received")
		default:
			serializeCommandBet, err := serializeCommandBet(bet)
			if err != nil {
				return err
			}
			binary.Write(bufferBetsData, binary.BigEndian, serializeCommandBet)
		}
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(BATCH_COMMAND))

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(bufferBetsData.Bytes()))); err != nil {
		return err
	}

	buffer.Write(bufferBetsData.Bytes())

	if len(buffer.Bytes()) > MAX_SIZE_BATCH {
		return fmt.Errorf("batch size is too big, max size is %d bytes", MAX_SIZE_BATCH)
	}

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol) ReceiveBatchResponse(amountBets int, exit chan os.Signal) (bool, error) {
	buf := make([]byte, 1)

	err := p.socket.Recvall(1, buf)
	if err != nil {
		return false, err
	}
	if buf[0] == RESPONSE_BATCH_COMMAND_OK {
		return true, nil
	}
	if buf[0] != RESPONSE_BATCH_COMMAND_ERROR {
		return false, fmt.Errorf("error processing batch")
	}
	return false, fmt.Errorf("invalid response command")
}

func (p *Protocol) SendDisconnect() error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(DISCONNECT_COMMAND))

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (p *Protocol) Close() error {
	return p.socket.Close()
}

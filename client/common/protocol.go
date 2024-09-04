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

	// comandos que envia
	BET_COMMAND          = 1
	BATCH_COMMAND        = 2
	DISCONNECT_COMMAND   = 3
	WAIT_WINNERS_COMMAND = 4
	CLIENT_ID            = 5

	// comandos que recibe
	RESPONSE_BET_COMMAND         = 1
	RESPONSE_BATCH_COMMAND_OK    = 2
	RESPONSE_BATCH_COMMAND_ERROR = 3
	RESPONSE_WINNERS_COMMAND     = 4
	RESPONSE_CLIENT_ID           = 5
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

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(serializeData))); err != nil {
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

func (p *Protocol) Close() error {
	return p.socket.Close()
}

func (p *Protocol) SendWaitingWinners() error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(WAIT_WINNERS_COMMAND))

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol) ReceiveWinners() ([]uint32, error) {
	buf := make([]byte, 1)

	err := p.socket.Recvall(1, buf)
	if err != nil {
		return nil, err
	}
	if buf[0] != RESPONSE_WINNERS_COMMAND {
		return nil, fmt.Errorf("invalid response command")
	}

	buf = make([]byte, 4)
	err = p.socket.Recvall(4, buf)
	if err != nil {
		return nil, err
	}

	dataSize := binary.BigEndian.Uint32(buf)

	buf = make([]byte, dataSize)
	err = p.socket.Recvall(int(dataSize), buf)
	if err != nil {
		return nil, err
	}

	documentsWinner := []uint32{}
	for i := 0; i < len(buf); i += 4 {
		document := binary.BigEndian.Uint32(buf[i : i+4])
		documentsWinner = append(documentsWinner, document)
	}

	return documentsWinner, nil
}

func (p *Protocol) SendId(id uint8) error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(CLIENT_ID))

	binary.Write(buffer, binary.BigEndian, id)

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

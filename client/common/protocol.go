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
	BATCH_COMMAND        = 2
	WAIT_WINNERS_COMMAND = 4
	CLIENT_ID            = 5

	// comandos que recibe
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

// SendBatch sends a batch of bets to the server.
// Encodings:
// - uint8: command
// - uint32: size of the batch
// - serialized bets
func (p *Protocol) SendBatch(bets []*Bet, exit chan os.Signal) error {
	bufferBetsData := new(bytes.Buffer)

	for _, bet := range bets {
		select {
		case <-exit:
			return fmt.Errorf("exit signal received")
		default:
			serializeData, err := bet.Serialize()
			if err != nil {
				return err
			}

			// Write the size of the bet data
			if err := binary.Write(bufferBetsData, binary.BigEndian, uint8(len(serializeData))); err != nil {
				return err
			}

			// Write the serialized bet data
			bufferBetsData.Write(serializeData)
		}
	}

	buffer := new(bytes.Buffer)

	// Write the command
	binary.Write(buffer, binary.BigEndian, uint8(BATCH_COMMAND))

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(bufferBetsData.Bytes()))); err != nil {
		return err
	}

	// Write the serialized bets
	buffer.Write(bufferBetsData.Bytes())

	// If the batch size is bigger than the maximum allowed, return an error
	if len(buffer.Bytes()) > MAX_SIZE_BATCH {
		return fmt.Errorf("batch size is too big, max size is %d bytes", MAX_SIZE_BATCH)
	}

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// ReceiveBatchResponse receives the response from the server after sending a batch.
// Returns true if the response is ok, and false otherwise.
func (p *Protocol) ReceiveBatchResponse(exit chan os.Signal) (bool, error) {
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

// SendWaitingWinners sends a command to the server to wait for the winners.
// Encodings:
// - uint8: command
func (p *Protocol) SendWaitingWinners() error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(WAIT_WINNERS_COMMAND))

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// ReceiveWinners receives the winners from the server.
// Return a slice with the document numbers of the winners.
func (p *Protocol) ReceiveWinners() ([]uint32, error) {
	buf := make([]byte, 1)

	// Receive the command
	err := p.socket.Recvall(1, buf)
	if err != nil {
		return nil, err
	}
	if buf[0] != RESPONSE_WINNERS_COMMAND {
		return nil, fmt.Errorf("invalid response command")
	}

	// Receive the size of the data
	buf = make([]byte, 4)
	err = p.socket.Recvall(4, buf)
	if err != nil {
		return nil, err
	}

	dataSize := binary.BigEndian.Uint32(buf)

	// Receive the data (documents)
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

// SendId sends the client id to the server.
// Encodings:
// - uint8: command
// - uint8: id
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

func (p *Protocol) Close() error {
	return p.socket.Close()
}
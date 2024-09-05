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

// SendDisconnect sends a disconnect command to the server.
func (p *Protocol) SendDisconnect() error {
	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, uint8(DISCONNECT_COMMAND))

	err := p.socket.Sendall(len(buffer.Bytes()), buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Close closes the connection to the server.
func (p *Protocol) Close() error {
	return p.socket.Close()
}

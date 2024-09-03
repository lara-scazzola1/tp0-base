package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Bet struct {
	Name       string
	Lastname   string
	Document   uint32
	BirthDay   uint8
	BirthMonth uint8
	BirthYear  uint16
	Number     uint32
	Agency     uint8
}

func NewBet(name string, lastname string, document uint32, birthDate string, number uint32, agency uint8) (*Bet, error) {
	date := strings.Split(birthDate, "-")
	if len(date) != 3 {
		return nil, fmt.Errorf("date must have the format YYYY-MM-DD")
	}

	year, err := strconv.ParseUint(date[0], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("error parsing year: %v", err)
	}
	month, err := strconv.ParseUint(date[1], 10, 8)
	if err != nil {
		return nil, fmt.Errorf("error parsing month: %v", err)
	}
	day, err := strconv.ParseUint(date[2], 10, 8)
	if err != nil {
		return nil, fmt.Errorf("error parsing day: %v", err)
	}

	return &Bet{
		Name:       name,
		Lastname:   lastname,
		Document:   document,
		BirthDay:   uint8(day),
		BirthMonth: uint8(month),
		BirthYear:  uint16(year),
		Number:     number,
		Agency:     agency,
	}, nil
}

func (b *Bet) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// Serializacion del nombre primero se manda la longitud de la cadena en 1 byte
	// y luego la cadena en si
	nameLength := uint8(len(b.Name))
	if err := binary.Write(buffer, binary.BigEndian, nameLength); err != nil {
		return nil, err
	}
	buffer.Write([]byte(b.Name))

	// Serializacion del apellido primero se manda la longitud de la cadena en 1 byte
	// y luego la cadena en si
	lastnameLength := uint8(len(b.Lastname))
	if err := binary.Write(buffer, binary.BigEndian, lastnameLength); err != nil {
		return nil, err
	}
	buffer.Write([]byte(b.Lastname))

	// Serializacion del documento en 4 bytes
	if err := binary.Write(buffer, binary.BigEndian, b.Document); err != nil {
		return nil, err
	}

	// Serializacion de la fecha de nacimiento en 4 bytes
	if err := binary.Write(buffer, binary.BigEndian, b.BirthDay); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, b.BirthMonth); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, b.BirthYear); err != nil {
		return nil, err
	}

	// Serializacion del numero de la apuesta en 4 bytes
	if err := binary.Write(buffer, binary.BigEndian, b.Number); err != nil {
		return nil, err
	}

	// Serializacion de la agencia en 1 byte
	if err := binary.Write(buffer, binary.BigEndian, b.Agency); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

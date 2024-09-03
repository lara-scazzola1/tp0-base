package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config   ClientConfig
	protocol *Protocol
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.protocol = NewProtocol(conn)
	return nil
}

// Gets the names of the files in a directory
func getFilenames(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	filenames := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			filenames = append(filenames, filepath.Join(dir, file.Name()))
		}
	}
	return filenames, nil
}

// Get the number of agency from the filename
func getAgency(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "-")
	return strings.TrimSuffix(parts[1], ".csv")
}

// Parse a line from a CSV file into a Bet struct
func parseLine(line string, agency string) (*Bet, error) {
	fields := strings.Split(line, ",")
	if len(fields) < 5 {
		return nil, fmt.Errorf("invalid line format: %v", line)
	}

	documento, err := strconv.ParseUint(fields[2], 10, 32)
	if err != nil {
		return nil, err
	}

	numero, err := strconv.ParseUint(fields[4], 10, 32)
	if err != nil {
		return nil, err
	}

	agencia, err := strconv.ParseUint(agency, 10, 8)
	if err != nil {
		return nil, err
	}

	bet, err := NewBet(
		fields[0],
		fields[1],
		uint32(documento),
		fields[3],
		uint32(numero),
		uint8(agencia),
	)
	if err != nil {
		return nil, err
	}

	return bet, nil
}

// Send a batch of bets to the server
func sendBatch(batch []*Bet, c *Client, exit chan os.Signal) error {
	if err := c.protocol.SendBatch(batch, exit); err != nil {
		return fmt.Errorf("error sending batch: %w", err)
	}

	ok, err := c.protocol.ReceiveBatchResponse(len(batch), exit)
	if err != nil {
		return fmt.Errorf("error receiving batch response: %w", err)
	}
	if !ok {
		log.Errorf("action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}")
		return nil
	}

	return nil
}

func processFile(file string, maxBatchSize int, c *Client, exit chan os.Signal) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	agency := getAgency(file)

	var batch []*Bet

	for scanner.Scan() {
		select {
		case <-exit:
			log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
			return nil
		default:
			bet, err := parseLine(scanner.Text(), agency)
			if err != nil {
				log.Errorf("error when parsing line: %v", err)
				continue
			}

			batch = append(batch, bet)

			if len(batch) >= maxBatchSize {
				if err := sendBatch(batch, c, exit); err != nil {
					return err
				}
				batch = []*Bet{}
				//time.Sleep(c.config.LoopPeriod)
			}
		}
	}

	if len(batch) > 0 {
		if err := sendBatch(batch, c, exit); err != nil {
			return fmt.Errorf("error sending the remaining batch: %w", err)
		}
	}

	return scanner.Err()
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(exit chan os.Signal, v *viper.Viper) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {
		case <-exit:
			log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
			return
		default:
			if err := c.createClientSocket(); err != nil {
				return
			}

			maxBatchSize := v.GetInt("batch.maxAmount")

			files, err := getFilenames("dataset")
			if err != nil {
				log.Errorf("Error getting filenames: %v", err)
				return
			}

			for _, file := range files {
				err := processFile(file, maxBatchSize, c, exit)
				if err != nil {
					log.Errorf("Error processing file %s: %v", file, err)
					return
				}
			}

			c.protocol.Close()
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

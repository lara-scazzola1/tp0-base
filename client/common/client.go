package common

import (
	"net"
	"os"
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
			// Create the connection the server in every loop iteration. Send an
			err := c.createClientSocket()
			if err != nil {
				return
			}

			bet, err := NewBet(
				v.GetString("nombre"),
				v.GetString("apellido"),
				v.GetUint32("documento"),
				v.GetString("nacimiento"),
				v.GetUint32("numero"),
				uint8(v.GetUint32("agencia")),
			)

			if err != nil {
				log.Errorf("action: create_bet | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			err = c.protocol.SendBet(bet)
			if err != nil {
				log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			ok, err := c.protocol.ReceiveResponseBet()
			if err != nil {
				log.Errorf("action: receive_response_bet | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}
			if !ok {
				log.Errorf("action: receive_response_bet | result: fail | client_id: %v | error: invalid_response",
					c.config.ID,
				)
				return
			}

			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", v.GetUint32("documento"), v.GetUint32("numero"))

			c.protocol.Close()

			// Wait a time between sending one message and the next one
			time.Sleep(c.config.LoopPeriod)
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

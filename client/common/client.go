package common

import (
	"github.com/op/go-logging"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const maxRetries = 3

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
	config ClientConfig
	conn   net.Conn
	bet    Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
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
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChannel:
		log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
		if c.conn != nil {
			c.conn.Close()
		}
		return
	default:
		bet := BetData(c.config.ID)
		if bet == nil {
			return
		}

		client_message := FormatMessage(*bet)
		if client_message == "" {
			return
		}

		var ack string
		var err error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := c.createClientSocket(); err != nil {
				log.Errorf("action: client_connect | result: fail | client_id: %v | attempt: %d/%d | error: %v",
					c.config.ID, attempt, maxRetries,
					err,
				)
				continue
			}

			if err := SendClientMessage(c.conn, client_message); err != nil {
				log.Errorf("action: send_message | result: fail | client_id: %v | attempt: %d/%d | error: %v",
					c.config.ID, attempt, maxRetries,
					err,
				)
				c.conn.Close()
				continue
			}

			ack, err = RecieveServerAck(c.conn)
			if err != nil {
				log.Errorf("action: receive_server_ack | result: fail | client_id: %v | attempt: %d/%d | error: %v", c.config.ID, attempt, maxRetries, err)
				c.conn.Close()
				time.Sleep(1 * time.Second) //sleep before next attempt
				continue
			}

			break
		}

		if err != nil {
			log.Criticalf("action: apuesta_enviada | result: fail | client_id: %v | error: could not send client message after %d retries",
				c.config.ID, maxRetries)
			if c.conn != nil {
				c.conn.Close()
			}
			return
		}

		if CheckServerAck(ack, c.bet) {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		} else {
			log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		}

		c.conn.Close()
	}
}

package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

const maxRetries = 3

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	bets   []Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bets []Bet) *Client {
	client := &Client{
		config: config,
		bets:   bets,
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

func (c *Client) StartClientLoop() {

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	agencyBets := len(c.bets)
	betCount := 0

	if err := c.createClientSocket(); err != nil {
		return
	}

	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for i := 0; i < agencyBets; i += c.config.BatchMaxAmount {
		select {
		case <-sigChannel:
			log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
			if c.conn != nil {
				c.conn.Close()
			}
			return
		default:
			end := i + c.config.BatchMaxAmount
			if end > agencyBets {
				end = agencyBets
			}

			batch := c.bets[i:end]
			betCount += len(batch)

			message := FormatBatchMessage(batch, betCount)

			if err := SendClientMessage(c.conn, message); err != nil {
				log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			ack, err := RecieveServerAck(c.conn)
			if err != nil {
				log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v",
					c.config.ID, err)
				return
			}

			success, batchSize := CheckBatchServerResponse(ack)
			if !success {
				log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | amount: %v",
					c.config.ID, batchSize)
				return
			} else {
				log.Infof("action: apuesta_enviada | result: success | client_id: %v | amount: %v",
					c.config.ID, batchSize)
			}

			time.Sleep(c.config.LoopPeriod)
		}

	}

	endOfFileMsg := FormatEndMessage(c.config.ID)
	if err := SendClientMessage(c.conn, endOfFileMsg); err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: could not send end of batch sending to server %v",
			c.config.ID,
			err,
		)
		return
	}

	ack, err := RecieveServerAck(c.conn)
	if err != nil {
		log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	if CheckEndServerResponse(ack, agencyBets, c.config.ID) {
		log.Infof("action: apuesta_enviada | result: success | client_id: %v | amount: %v", c.config.ID, agencyBets)
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | amount: %v", c.config.ID, agencyBets)
	}

}

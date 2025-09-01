package common

import (
	"fmt"
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

	var err error = nil

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	if err := c.createClientSocket(); err != nil {
		return
	}

	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	message := FormatBetSendingMessage()

	if err := SendClientMessage(c.conn, message); err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	err = c.SendClientBets(sigChannel)
	if err != nil {
		return
	}

	err = c.SendEndOfBetsMessage(sigChannel)
	if err != nil {
		return
	}

	ack, err := RecieveServerAck(c.conn)
	if err != nil {
		log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	agencyBets := len(c.bets)

	if CheckEndServerResponse(ack, agencyBets, c.config.ID) {
		log.Infof("action: apuesta_enviada | result: success | client_id: %v | amount: %v", c.config.ID, agencyBets)
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | amount: %v", c.config.ID, agencyBets)
		return
	}

	c.conn.Close()

	err = c.WaitForLoteryResults(sigChannel)
	if err != nil {
		return
	}

}

func (c *Client) SendClientBets(sigChannel chan os.Signal) error {
	agencyBets := len(c.bets)
	betCount := 0

	for i := 0; i < agencyBets; i += c.config.BatchMaxAmount {
		select {
		case <-sigChannel:
			log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
			return fmt.Errorf("client_shutdown") //es para que entre en el defer cierre y salga
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
				return fmt.Errorf("could not send bet batch: %w", err)
			}

			ack, err := RecieveServerAck(c.conn)
			if err != nil {
				log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v",
					c.config.ID, err)
				return fmt.Errorf("receive server ack failed for client batch sending: %w", err)
			}

			success, batchSize := CheckBatchServerResponse(ack)
			if !success {
				log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | amount: %v",
					c.config.ID, batchSize)
				return fmt.Errorf("batch failed at size %v", batchSize)
			} else {
				log.Infof("action: apuesta_enviada | result: success | client_id: %v | amount: %v",
					c.config.ID, batchSize)
			}

			//time.Sleep(c.config.LoopPeriod)

		}

	}
	return nil
}

func (c *Client) SendEndOfBetsMessage(sigChannel chan os.Signal) error {
	select {
	case <-sigChannel:
		log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
		return fmt.Errorf("client_shutdown")
	default:
		endOfFileMsg := FormatEndMessage(c.config.ID)
		if err := SendClientMessage(c.conn, endOfFileMsg); err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: could not send end of batch sending to server %v",
				c.config.ID,
				err,
			)
			return fmt.Errorf("error sending end of bets file message")
		}
	}

	return nil
}

func (c *Client) WaitForLoteryResults(sigChannel chan os.Signal) error {
	sleepTimer := 200 * time.Millisecond
	for {
		select {
		case <-sigChannel:
			log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
			return fmt.Errorf("client_shutdown")
		default:
			log.Infof("action: consulta_ganadores | result: in_progress | client_id: %v", c.config.ID)

			if c.conn != nil {
				c.conn.Close()
			}
			if err := c.createClientSocket(); err != nil {
				return err
			}

			loteryWinnerConsult := FormatWinnerConsult(c.config.ID)
			if err := SendClientMessage(c.conn, loteryWinnerConsult); err != nil {
				log.Errorf("action: send_message | result: fail | client_id: %v | error: could not send lotery winner consult to server %v",
					c.config.ID,
					err,
				)
				return fmt.Errorf("error sending end of bets file message")
			}

			//aca tiene que ir un wait for server, y el server le responde
			//con un error o con un result. Segun eso logguea al winner cierra la conexio
			//y sale (break), o si devuelve error, cierra la conexion mete un sleep y reintetna
			ack, err := RecieveServerAck(c.conn)
			if err != nil {
				log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v",
					c.config.ID, err)
				return fmt.Errorf("receive server ack for lottery response failed: %w", err)
			}
			success, winners := CheckLotteryResult(ack)
			if !success {
				log.Infof("action: consulta_ganadores | result: fail | status: not ready")
				c.conn.Close()
				time.Sleep(sleepTimer)
				sleepTimer *= 2

			} else {
				log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))
				c.conn.Close()
				return nil
			}
		}
	}
}

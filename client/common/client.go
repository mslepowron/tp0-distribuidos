package common

import (
	//"bufio"
	//"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
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
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
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

			// Create the connection the server in every loop iteration. Send an
			c.createClientSocket()

			// TODO: Modify the send to avoid short-write
			// fmt.Fprintf(
			// 	c.conn,
			// 	"[CLIENT %v] Message N°%v\n",
			// 	c.config.ID,
			// 	msgID,
			// ) //en esta lectura tengo que asegurarme que tampoco haya short read de la rta del server
			// msg, err := bufio.NewReader(c.conn).ReadString('\n')

			if err := SendClientMessage(c.conn, client_message); err != nil {
				log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				c.conn.Close()
				return
			}

			ack, err := RecieveServerAck(c.conn)
			if err != nil {
				log.Errorf("action: receive_server_ack | result: fail | client_id: %v | error: %v", c.config.ID, err)
				c.conn.Close()
				return
			}

			//aca habria que formatear el mensaje de bet, borrar este fmt.Fprintf
			//despues llamar a una funcion de send que gestione el short read
			//con un prefijo de cantidad de bytes o algo asi.

			if err != nil {
				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}

			println("Client %v received ack: %s\n", c.config.ID, ack) //solo pongo el print para que no tire error. Hay que decodear el ack del server

			c.conn.Close()

			//EXTRAR CAMPOS ACK DEL SERVER: DOCUMENTO Y NUMERO para imprimir log de succes o fail

			// Wait a time between sending one message and the next one
			time.Sleep(c.config.LoopPeriod)
		}

	}
}

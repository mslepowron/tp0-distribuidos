package common

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const MaxMessageSize = 8192

type Bet struct {
	AgencyId  string
	Name      string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

// BetData creates a Bet struct with the client's bet information stored in environment variables
func BetData(clientID string) *Bet {

	betName := os.Getenv("NOMBRE")
	betLastName := os.Getenv("APELLIDO")
	betDocument := os.Getenv("DOCUMENTO")
	betBirthDate := os.Getenv("NACIMIENTO")
	betNumber := os.Getenv("NUMERO")

	if betName == "" || betLastName == "" || betDocument == "" || betBirthDate == "" || betNumber == "" {
		log.Critical("Faltan variables de entorno requeridas para la apuesta")
		return nil
	}

	bet := &Bet{
		AgencyId:  clientID,
		Name:      betName,
		LastName:  betLastName,
		Document:  betDocument,
		BirthDate: betBirthDate,
		Number:    betNumber,
	}

	return bet
}

// FormatMessage formats the client's bet data into a protocol the server understands
func FormatMessage(bet Bet) string {

	for _, field := range []string{bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number} {
		if field == "" || strings.Contains(field, ";") {
			log.Critical("action: client_message_parser | result: fail | bet field cannot contain ';'")
			return ""
		}
		if field == "" {
			log.Critical("action: client_message_parser | result: fail | bet field but be complete")
			return ""
		}
	}

	msg := fmt.Sprintf("%s;%s;%s;%s;%s;%s", bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number)

	return msg
}

// WriteFull send all the data through the server connection while avoiding short-writes
func WriteFull(connection net.Conn, data []byte) error {
	totalWritten := 0
	dataLength := len(data)

	for totalWritten < dataLength {
		bytesWritten, err := connection.Write(data[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += bytesWritten
	}

	return nil
}

// SendClientMessage serializes the client data and sends it to the server
func SendClientMessage(connection net.Conn, message string) error {

	if len(message) > MaxMessageSize {
		log.Critical("action: send_client_message | result: fail | message is bigger than 8KB")
		return fmt.Errorf("message is bigger than 8KB")
	}

	messageBytes := []byte(message)

	messageLength := make([]byte, 4)
	binary.BigEndian.PutUint32(messageLength, uint32(len(messageBytes)))

	//try send message length
	if err := WriteFull(connection, messageLength); err != nil {
		return err
	}

	//send agency data
	return WriteFull(connection, messageBytes)
}

// RecieveServerAck read from the server connection de server response (ack) given to the
// message sent by the client
func RecieveServerAck(connection net.Conn) (string, error) {
	reader := bufio.NewReader(connection)

	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Critical("action: receive_server_ack | result: fail | error: reading Ack message")
		return "", fmt.Errorf("error reading Ack message: %w", err)
	}

	msg = strings.TrimSpace(msg)

	if len(msg) > 8192 {
		log.Critical("action: receive_server_ack | result: fail | server ack message is too large")
		return "", fmt.Errorf("server ack message is too large")
	}

	return msg, nil
}

// CheckServerAck compares the ack received from the server with the bet data to see if the information
// was recieved correctly.
func CheckServerAck(ack string, bet Bet) bool {
	serverAck := strings.Split(ack, ";")

	if len(serverAck) != 2 {
		log.Errorf("action: check_server_ack | result: fail | invalid ack format")
		return false
	}

	ackDocument, _ := strconv.Atoi(strings.TrimSpace(serverAck[0]))
	ackNumber, _ := strconv.Atoi(strings.TrimSpace(serverAck[1]))

	clientDocument, _ := strconv.Atoi(bet.Document)
	clientNumber, _ := strconv.Atoi(bet.Number)

	return ackDocument == clientDocument && ackNumber == clientNumber
}

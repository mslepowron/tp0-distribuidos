package common

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
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

// ReadAgencyBets reads the bets from a CSV file and returns a slice of Bet structs
func ReadAgencyBets(agencyId string) ([]Bet, error) {
	file, err := os.Open("agency.csv") //lo saca de la config del docker compose (volume)

	if err != nil {
		return nil, fmt.Errorf("action: read_agency_file | result: fail | error: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("action: read_agency_file | result: fail | error: reading CSV: %v", err)
	}

	var bets []Bet
	for _, line := range lines {
		if len(line) != 5 {
			return nil, fmt.Errorf("action: read_agency_file | result: fail | error: invalid csv format at line %v", line)
		}

		bet := Bet{
			AgencyId:  agencyId,
			Name:      line[0],
			LastName:  line[1],
			Document:  line[2],
			BirthDate: line[3],
			Number:    line[4],
		}
		bets = append(bets, bet)
	}

	return bets, nil
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

// FormatBatchMessage formats an agency bets data in a range of N bets, according
// to de batch size. It uses the same format as FormatMessage (established protocol with
// the server) and uses a delimiter to separate different bets.
func FormatBatchMessage(bets []Bet, betCount int) string {
	bets_string := make([]string, 0, len(bets))

	for _, bet := range bets {
		bet_message := FormatMessage(bet)
		bets_string = append(bets_string, bet_message)
	}

	agency_bets_message := strings.Join(bets_string, "\n")

	return agency_bets_message
}

func FormatBatches(bets []Bet, maxBatchSize int) []string{
	var messages []string
    var currentBatch []string
    currentSize := 0

	for _, bet := range bets {
        betMsg := FormatMessage(bet)
        betSize := len(betMsg) + 1

        
        if len(currentBatch) > 0 && 
           (len(currentBatch) >= maxBatchSize || currentSize+betSize > MaxMessageSize) {
            messages = append(messages, strings.Join(currentBatch, "\n"))
            currentBatch = []string{}
            currentSize = 0
        }

        currentBatch = append(currentBatch, betMsg)
        currentSize += betSize
    }

    if len(currentBatch) > 0 {
        messages = append(messages, strings.Join(currentBatch, "\n"))
    }

    return messages
}

func FormatEndMessage(agencyId string) string {
	return fmt.Sprintf("END_OF_FILE;%s", agencyId)
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

// CheckEndServerResponse checks the server response to the END_OF_FILE message
func CheckEndServerResponse(ack string, betCount int, agencyId string) bool {
	serverAck := strings.Split(ack, ";")

	if len(serverAck) != 2 {
		log.Errorf("action: check_server_ack | result: fail | invalid ack format")
		return false
	}

	ackBetsAmount, _ := strconv.Atoi(strings.TrimSpace(serverAck[0]))
	ackAgencyID, _ := strconv.Atoi(strings.TrimSpace(serverAck[1]))

	agencyIDInt, _ := strconv.Atoi(agencyId)

	return ackBetsAmount == betCount && ackAgencyID == agencyIDInt
}

// CheckBatchServerResponse checks the server response to a batch message
func CheckBatchServerResponse(ack string) (success bool, batchSize int) {

	ack = strings.TrimSpace(ack)
	serverAck := strings.Split(ack, ";")
	if strings.HasPrefix(ack, "BATCH_OK") {
		success = true
		if len(serverAck) == 2 {
			batchSize, _ = strconv.Atoi(strings.TrimSpace(serverAck[1]))
		}
	} else if strings.HasPrefix(ack, "ERROR_BATCH") {
		success = false
		if len(serverAck) == 2 {
			batchSize, _ = strconv.Atoi(strings.TrimSpace(serverAck[1]))
		}
	} else {
		// EOF o cualquier otro caso
		success = true // o false según lo que quieras manejar
		batchSize = 0
	}
	return
}

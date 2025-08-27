package common

import (
	"fmt"
	"os"
	"strings"
	// "github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

type Bet struct {
	AgencyId  string
	Name      string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

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

	log.Criticalf(
		"action: connect | result: fail | missing bet information for agency: %v", clientID,
	)

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

//Send bet info:
//serialize bet data -> send to server -> wait for confirmation (log)

func FormatMessage(bet Bet) string {

	for _, field := range []string{bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number} {
		if field == "" || strings.Contains(field, "|") {
			log.Critical("action: client_message_parser | result: fail | bet field cannot contain '|'")
			return ""
		}
		if field == "" {
			log.Critical("action: client_message_parser | result: fail | bet field but be complete")
			return ""
		}
	}

	msg := fmt.Sprintf("%s|%s|%s|%s|%s|%s", bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number)

	return msg
}

package main

import (
	"strings"
)

//********************************************Adding commands********************************************
func listRunCommands(author string, conversationID string, params ...string) {

	name := strings.Fields(author)

	message := "Yo " + name[0] + "\n \n <b>Try following commands</b> \n\n yumi run buildium \n\n yumi run password"

	addNote(conversationID, message)
}

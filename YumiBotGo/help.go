package main

import (
	"strings"
)

//********************************************Adding commands********************************************
func listRunCommands(author string, conversationID string, params ...string) {

	name := strings.Fields(author)

	message := "Yo " + name[0] + "\n \n <b>Try following commands</b>"
	for cmd, _ := range commands {
		message += "\n\n yumi run " + cmd
	}

	addNote(conversationID, message)
}

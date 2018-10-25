package main

import (
	"math/rand"
	"strings"
)

//********************************************Adding commands********************************************
func listRunCommands(author string, conversationID string, params ...string) {
	name := strings.Fields(author)
	var salution = []string{"Howdy", "Namaste", "Yo", "Aloha", "Ni Hau", "Star lord", "Happy", "Bonjour", "Hola", "Hallo"}

	message := "<h2>" + salution[rand.Intn(len(salution))] + " " + name[0] + "</h2>\n \n <b>Try following commands</b>"
	message += "\n\n reply-reply to a customer \n get- add an internal/team note"
	message += "\n\n <i>Arguments are optional</i>"

	for cmd, detail := range replyCommands {
		message += "<b>\n\n yumi reply " + cmd + "</b>\n\n" + detail.Description
	}

	for cmd, detail := range getCommands {
		message += "<b>\n\n yumi get " + cmd + "</b>\n\n" + detail.Description
	}
	message += "\n\n"
	message += "<a href=\"https://hpy.io/yumi\">Feedback/Report incorrect information</a>"

	addNote(conversationID, message)
}

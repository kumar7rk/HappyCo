package main

import (
	"strings"
)

//********************************************Adding commands********************************************
func listRunCommands(author string, conversationID string, params ...string) {
	name := strings.Fields(author)

	message := "Yo! " + name[0] + "\n \n <b>Try following commands</b>"
	message += "\n\n rep-reply to customer \n get- add a note"
	for cmd, detail := range repCommands {
		message += "<b>\n\n yumi rep " + cmd + "</b>\n\n" + detail.Description
	}

	for cmd, detail := range getCommands {
		message += "<b>\n\n yumi get " + cmd + "</b>\n\n" + detail.Description
	}
	message += "\n\n"
	message += "<a href=\"https://hpy.io/yumi\">Feedback/Report incorrect information</a>"

	addNote(conversationID, message)
}

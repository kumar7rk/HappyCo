package main

import (
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
	"os"
	"strings"
)
//********************************************Adding commands********************************************
func listRunCommands(author string, conversationID string, params ...string) {
	
	name :=strings.Fields(author)
	
	message := "Yo " + name[0] + "\n \n <b>Try following commands</b> \n\n yumi run buildium \n\n yumi run password"

	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_NOTE, message)
	
	if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
		fmt.Fprintf(os.Stderr, "Error from Intercom when running yumi help%v: %v\n", "", err)
		return
	}

}

package main

import (
	"fmt"
	intercom "gopkg.in/intercom/intercom-go.v2"
	"os"
	"strings"
)

//********************************************Sending password reply********************************************

func sendPasswordReply(user User, conversationID string, params ...string) {
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic := intercom.NewClient(accessToken, "")
	// if you don't give a name use the name user have in his account

	if len(params) == 0 {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	if len(params) > 0 {
		name = params[0]
	}
	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='http://intercom.help/happyco/frequently-asked-questions/happy-inspector-faq-setup-and-user-management/faq-managing-passwords-user-details-and-access-to-your-properties'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to speak to a human....... just reply</i>"

	_, err := ic.Conversations.Reply(conversationID, intercom.Admin{ID: "207278"}, intercom.CONVERSATION_COMMENT, message)

	if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
		fmt.Fprintf(os.Stderr, "Error from Intercom when replying to Buildium %v: %v\n", "", err)
		return
	}

}

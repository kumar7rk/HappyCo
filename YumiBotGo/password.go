package main

import (
	intercom "gopkg.in/intercom/intercom-go.v2"
	"strings"
)

//********************************************Sending password reply********************************************

func init() {
	commands["password"] = sendPasswordReply
}

func sendPasswordReply(user User, conversationID string, params ...string) {
	// if you don't give a name use the name user have in his account
	if len(params) == 0 && user.Email != "" {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	if len(params) > 0 {
		name = params[0]
	}
	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='http://intercom.help/happyco/frequently-asked-questions/happy-inspector-faq-setup-and-user-management/faq-managing-passwords-user-details-and-access-to-your-properties'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to speak to a human....... just reply</i>"

	addReply(conversationID, message)
	ic.Conversations.Assign(conversationID, &intercom.Admin{ID: "207278"}, &intercom.Admin{ID: "931140"})
}

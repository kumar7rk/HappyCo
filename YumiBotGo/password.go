package main

import (
	"strings"
)

//********************************************Init********************************************
func init() {
	commands["password"] = sendPasswordReply
}

//********************************************Sending password reply********************************************
func sendPasswordReply(user User, conversationID string, params ...string) {
	// if you don't give a name use the name user have in his account
	if len(params) == 0 && user.Email != "" {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	if len(params) > 0 {
		name = params[0]
	}
	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to contact a human....... just reply</i>"

	addReply(conversationID, message)
	assignConversation(conversationID, "1398520")
	snoozeConversation(conversationID)
}

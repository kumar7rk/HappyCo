package main

import (
	"strconv"
	"strings"
)

//********************************************Init********************************************
func init() {
	repCommands["password"] = Command{Func: sendPasswordReply, Description: `Sends password reset instructions.
Argument - same as buildium`}
}

//********************************************Sending password reply********************************************
func sendPasswordReply(user User, conversationID string, params ...string) {
	// if you don't give a name use the name user have in his account
	if len(params) == 0 && user.Email != "" {
		params = strings.Split(user.Name, " ")
	}
	name := "there"
	snoozeDays := int64(3)

	if len(params) > 0 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = params[0]
		} else {
			name = strings.Split(user.Name, " ")[0]
			snoozeDays = int64(val)
		}
	}
	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to contact a human....... just reply</i>"

	addReply(conversationID, message)
	assignConversation(conversationID, "1398520")
	snoozeConversation(conversationID, snoozeDays)
}
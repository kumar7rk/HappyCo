package main

import (
	"fmt"
	"strconv"
	"strings"
)

//********************************************Init********************************************
func init() {
	repCommands["password"] = Command{Func: sendPasswordReply, Description: `Sends password reset instructions.
		<b>Default</b> snooze is 3 days.
		<b>Arguments</b> either name or snooze time.
		<b>Example</b> yumi rep password [name] or [days to snooze]`}
}

//********************************************Sending password reply********************************************
func sendPasswordReply(user User, conversationID string, params ...string) {
	name := "there"
	snoozeDays := int64(3)
	if user.Email != "" {
		name = strings.Split(user.Name, " ")[0]
	}
	if len(params) == 2 {
		//assuming first param is a name
		name = params[0]
		//checking if the second param is a number or not
		val, err := strconv.Atoi(params[1])
		if err != nil {
			fmt.Printf("Wrong parameters: %v\n", err)
		} else {
			snoozeDays = int64(val)
		}
	}
	if len(params) == 1 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = params[0]
		} else {
			snoozeDays = int64(val)
		}
	}

	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to contact a human....... just reply</i>"

	addReply(conversationID, message)
	assignConversation(conversationID, "1398520")
	snoozeConversation(conversationID, snoozeDays)
}

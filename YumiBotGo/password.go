package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const P3BoxID = "1398520"

//********************************************Init********************************************
func init() {
	repCommands["password"] = RepCommand{Func: sendPasswordReply, Description: `Sends password reset instructions.
		<b>Default</b> snooze is 3 days.
		<b>Arguments</b> either name or snooze time.
		<b>Example</b> yumi rep password [name] or [days to snooze]`}
}

//********************************************Sending password reply********************************************

func sendPasswordReply(user User, author Author, conversationID string, params ...string) {
	name := "there"
	snoozeDuration := 3 * 24 * time.Hour

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
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}
	if len(params) == 1 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = params[0]
		} else {
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}

	message := passwordMessage(name)

	addReply(conversationID, YumiBot.ID, message)
	assignConversation(conversationID, P3BoxID)
	snoozeConversation(conversationID, snoozeDuration)
}

package main

import (
	"strconv"
	"strings"
	"time"

	"happyco/libs/log"
)

// P2ABoxID represents intercom's P2A inbox ID
const P2ABoxID = "1615207"

//********************************************Init********************************************
func init() {
	repCommands["buildium"] = RepCommand{Func: sendBuildiumReply, Description: `Sends a buildium message to a user
		<b>Default</b> snooze is 7 days.
		<b>Arguments</b> either name or snooze time.
		<b>Example</b> yumi rep buildium [name] or [days to snooze]`}
}

//********************************************Sending Buildium reply********************************************
func sendBuildiumReply(user User, author Author, conversationID string, params ...string) {
	name := "there"
	snoozeDuration := 7 * 24 * time.Hour

	if hasIdentifiableName(user) {
		name = strings.Split(user.Name, " ")[0]
	}
	if len(params) == 2 {
		//assuming first param is a name
		name = params[0]
		//checking if the second param is a number or not
		val, err := strconv.Atoi(params[1])
		if err != nil {
			log.Error.KV("err",err).KV("params", params).KV("conversationID", conversationID).Println("could not parse number of snooze days for sending buildium reply")
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

	addReply(yumiBot.ID, conversationID, buildiumMessage(name))
	assignConversation(conversationID, P2ABoxID)
	snoozeConversation(conversationID, snoozeDuration)
}

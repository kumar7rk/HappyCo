package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//********************************************Init********************************************
func init() {
	repCommands["buildium"] = Command{Func: sendBuildiumReply, Description: `Sends a buildium message to a user
	 <b>Default</b> snooze is 7 days.
	 <b>Arguments</b> either name or snooze time.
	 <b>Example</b> yumi rep buildium [name] or [days to snooze]`}
}

//********************************************Sending Buildium reply********************************************
func sendBuildiumReply(user User, conversationID string, params ...string) {
	name := "there"
	snoozeDuration := 7 * 24 * time.Hour

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

	message := "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"

	addReply(conversationID, message)
	assignConversation(conversationID, "1615207")
	snoozeConversation(conversationID, snoozeDuration)
}

package main

import (
	"strconv"
	"strings"
	"time"

	"happyco/libs/log"
)

// P3BoxID represents intercom's P3 inbox ID
const P3BoxID = "1398520"

//********************************************Init********************************************
func init() {
	replyCommands["password"] = ReplyCommand{Func: sendPasswordReply, Description: `Sends password reset instructions.
		<b>Default</b> snooze is 3 days.
		<b>Arguments</b> either name or snooze time.
		<b>Example</b> yumi reply password [name] or [days to snooze]`}
}

//********************************************Sending password reply********************************************
func sendPasswordReply(user User, author Author, conversationID string, params ...string) {
	name := "there"
	snoozeDuration := 3 * 24 * time.Hour

	if hasIdentifiableName(user) {
		name = strings.Split(user.Name, " ")[0]
	}
	if len(params) == 2 {
		//assuming first param is a name
		name = strings.Title(params[0])
		//checking if the second param is a number or not
		val, err := strconv.Atoi(params[1])
		if err != nil {
			log.Error.KV("err", err).KV("params", params).KV("conversationID", conversationID).Println("could not parse number of snooze days for sending password reply")
		} else {
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}
	if len(params) == 1 {
		val, err := strconv.Atoi(params[0])
		if err != nil {
			name = strings.Title(params[0])
		} else {
			snoozeDuration = time.Duration(val) * 24 * time.Hour
		}
	}

	message := passwordMessage(name)

	addReply(yumiBot.ID, conversationID, message)
	assignConversation(conversationID, P3BoxID)
	snoozeConversation(conversationID, snoozeDuration)
}

func hasIdentifiableName(user User) bool {

	firstName := strings.Split(user.Name, " ")[0]
	userNameInLowerCase := strings.ToLower(user.Name)

	//checking certain keywords and numbers
	var excludeList = []string{"property", "inspector", "management", "maintenance", "department", "mgmt", "dept", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, name := range excludeList {
		if strings.Contains(userNameInLowerCase, name) {
			return false
		}
	}

	//checking acronyms
	if firstName == strings.ToUpper(firstName) {
		return false
	}
	if user.Email == "" {
		return false
	}
	return true
}

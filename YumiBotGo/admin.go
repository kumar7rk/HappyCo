package main

//********************************************Init********************************************
func init() {
	getCommands["admin"] = Command{Func: showAllAdmins, Description: `Prints a note with the name and emails of all the admins in a business.
		<b>Arguments</b> userID.
		<b>Example</b> yumi get admin [userID]`}
}

//********************************************Adding commands********************************************
func showAllAdmins(user User, conversationID string, params ...string) {
	var message string
	userID := user.UserID
	if len(params) > 0 {
		userID = params[0]
	}
	adminsRec := getAdmins(userID)
	for _, admin := range adminsRec {
		message += admin.Detail
		message += "\n"
	}
	addNote(conversationID, message)
}

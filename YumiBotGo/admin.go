package main

//********************************************Init********************************************
func init() {
	getCommands["admin"] = Command{Func: showAllAdmins, Description: `Prints a note with the name and emails of all the admins in a business.`}
}

//********************************************Adding commands********************************************
func showAllAdmins(user User, conversationID string, params ...string) {
	var message string
	adminsRec := getAdmins(user.UserID)
	for _, admin := range adminsRec {
		message += admin.Detail
		message += "\n"
	}
	addNote(conversationID, message)
}

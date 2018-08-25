package main

import (
	"time"
)
func init() {
	commands["report"] = showRecentReports
}
//********************************************Adding commands********************************************
func showRecentReports(user User, conversationID string, params ...string) {

	reportsRec := getReports(user.UserID)
	message := "<b>Showing recent reports (max 5) created by "+ user.Name+" in last 30 days</b>"
	message += "\n"
	for _, report := range reportsRec {
		var url = "https://manage.happyco.com/reports/" + report.PublicID
		var date, _ = time.Parse(time.RFC3339, report.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		message += "<a href=" + url + ">" + report.Name + "</a>" + " " + formattedDate
		message += "\n"
	}
	addNote(conversationID, message)
}
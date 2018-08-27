package main

import (
	"strconv"
	"time"
)

func init() {
	commands["report"] = showRecentReports
}

//********************************************Adding commands********************************************
func showRecentReports(user User, conversationID string, params ...string) {
	var limit = 5
	if len(params) > 0 {
		limit, _ = strconv.Atoi(params[0])
	}
	reportsRec := getReports(user.UserID, limit)
	message := "<b>Showing " + strconv.Itoa(limit) + " recent reports created by " + user.Name + " in last 30 days</b>"
	message += "\n"
	for _, report := range reportsRec {
		var url = "https://manage.happyco.com/reports/" + report.PublicID
		var date, _ = time.Parse(time.RFC3339, report.CreatedAt)
		var formattedDate = date.Format("02 Jan 2006 3:04PM")

		message += "<a href=\"" + url + "\">" + report.Name + "</a>" + " " + formattedDate
		message += "\n"
	}
	addNote(conversationID, message)
}

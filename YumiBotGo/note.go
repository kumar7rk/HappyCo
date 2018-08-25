package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"time"
	// "os"
)



//********************************************Adding note to conversation********************************************

func init() {
	commands["note"] = makeAndSendNote
}

func makeAndSendNote(user User, conversationID string, params ...string) {
	fmt.Println("makeAndSendNote")
	ID := user.UserID
	p := fmt.Println

	//testing prints
	p("Conversation id: " + conversationID)
	p("User id: " + ID)
	p("User name: " + user.Name)
	p("User email: " + user.Email)

	// calling the method to compile the note with all the required information
	note, _ := makeNote(ID)

	addNote(conversationID, note)
}

//********************************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(ID string) (inspectionsRec []Inspection, reportsRec []Report, businessRec []Business, iapRec []IAP, integrationName string, planTypeRec []Plan) {
	fmt.Println("getUserData")
	//fetching most recent (5) inspections for the user within the last 30 days.
	inspectionsRec = getInspections(ID)

	//fetching most recent (5) reports for the user within the last 30 days.
	reportsRec = getReports(ID)

	// fetching business id and role id for user role in this business
	businessRec = getBusiness(ID)

	// checking if the business is on IAP get the expiry date
	iapRec = getIAP(ID)

	// Check if the business has integration w/Yardi
	integrationName = getIntegration(ID)	

	// Plan type = DD/Buildium/MRI
	planTypeRec = getUserPlanType(ID)
	return
}

//********************************************Building note********************************************

// build the note in a string format
func makeNote(us_id string) (string, string) {
	fmt.Println("makeNote")
	var note string
	var formattedDate string

	//getting user data from the database
	inspectionsRec, reportsRec, businessRec, iapRec, integrationName, planTypeRec := getUserData(us_id)

	//******************constructing business string******************
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"

	note += "<b>✅User is associated with the following businesses</b><br/><br/>"

	// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
	// 8, 9 = Basic Admin, PM
	var roles = []string{"", "Admin", "Process Manger", "Inspector", "Limited Inspector", "", "", "", "Admin", "Process Manager"}

	for _, business := range businessRec {

		permission, err := strconv.Atoi(business.Role)
		if err != nil {
			panic(err)
		}
		var roleID = business.Role

		if roleID == "1" || roleID == "2" || roleID == "3" || roleID == "4" {
			note += " The business is on Constant/full Permissions \n"
		}

		if roleID == "8" || roleID == "9" {
			note += " The business is on Basic Permissions \n"
		}

		note += roles[permission] + " for "
		var text = "https://manage.happyco.com/admin/businesses/" + business.ID
		note += text + "<br/><br/>"
	}
	note += "\n"

	//******************constructing integration string******************
	if integrationName != "" {
		note += "\n"
		note += "\n"
		note += "The business is " + integrationName
	}

	//******************constructing plan type string******************

	planType := "plan type"
	for _, plan := range planTypeRec {
		if plan.Type == "due_diligence" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Due Diligence"
		}
		if plan.Type == "buildium" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Buildium"
			planType = plan.Type
		}
		if plan.Type == "mri" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "MRI"
		}
	}

	note += "\n"
	note += "\n"

	//******************constructing inspection string******************
	note += "<b>✅   Yumi found these recent (max: 5) <em>Inspections in last 30 days:</em></b><br/>"
	note += "\n"
	var url string
	for _, inspection := range inspectionsRec {

		url = "https://manage.happyco.com/folder/" + inspection.FolderID + "/inspections/" + inspection.ID

		var date, _ = time.Parse(time.RFC3339, inspection.CreatedAt)
		formattedDate = date.Format("02 Jan 2006 3:04PM")

		note += "<a href=" + url + ">" + url + "</a>" + " " + formattedDate
		note += "\n"
	}
	//******************constructing report string******************

	note += "\n"
	note += "\n"
	note += "<b>✅   Yumi found these recent (max: 5) <em>Reports in last 30 days:</em></b><br/>"
	note += "\n"
	for _, report := range reportsRec {

		var url = "https://manage.happyco.com/reports/" + report.PublicID

		var date, _ = time.Parse(time.RFC3339, report.CreatedAt)
		formattedDate = date.Format("02 Jan 2006 3:04PM")

		note += "<a href=" + url + ">" + report.Name + "</a>" + " " + formattedDate
		note += "\n"
	}

	//******************constructing iap string******************

	for _, iap := range iapRec {
		if iap.Expiry != "" {
			var date, _ = time.Parse(time.RFC3339, iap.Expiry)
			formattedDate = date.Format("02 Jan 2006 3:04PM")
			note += "\n"
			note += "\n"
			note += "<b>The business is on IAP. It expires on </b>" + formattedDate
		}
	}
	note += "\n"
	note += "<a href=" + "https://hpy.io/yumi" + ">" + "Feedback/Report incorrect information" + "</a>"
	return note, planType
}
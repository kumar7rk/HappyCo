package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"time"
	"strings"
)



//********************************************Adding note to conversation********************************************

func init() {
	commands["note"] = makeAndSendNote
}

func makeAndSendNote(user User, conversationID string, params ...string) {
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
	//fetching most recent (5) inspections for the user within the last 30 days.
	inspectionsRec = getInspections(ID,5)

	//fetching most recent (5) reports for the user within the last 30 days.
	reportsRec = getReports(ID,5)

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
	var note string
	var formattedDate string

	//getting user data from the database
	_, _, businessRec, iapRec, integrationName, planTypeRec := getUserData(us_id)

	//******************constructing business string******************
	note = "<b>🐶Note</b><br/><br/>"

	// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
	// 8, 9 = Basic Admin, PM
	var roles = []string{"", "Admin", "Process Manger", "Inspector", "Limited Inspector", "", "", "", "Admin", "Process Manager"}

	for _, business := range businessRec {

		permission, err := strconv.Atoi(business.Role)
		if err != nil {
			panic(err)
		}
		var roleID = business.Role
		var BusinessPermission string
		if roleID == "1" || roleID == "2" || roleID == "3" || roleID == "4" {
			BusinessPermission = "Contant/Full"
		}

		if roleID == "8" || roleID == "9" {
			BusinessPermission = "Basic"
		}
		note+="<b>Business: </b>" + business.Name +"\n"
		note+="<b>BusinessID:</b>" + business.ID +"\n"
		note+="<b>Permissions:</b>" + BusinessPermission +"\n"
		note+="<b>Role:</b>" + roles[permission] +"\n"
	}	
	//******************constructing plan type string******************

	planType := "plan type"
	for _, plan := range planTypeRec {
		if plan.Type == "buildium" {
			planType = plan.Type
		}
		plan.Type = strings.Replace(plan.Type,"_"," ",-1)
		note += "<b>Plan: </b>" + strings.Title(plan.Type) +"\n"
	}
	
	note+="<b>MRR:</b>" + "None" +"\n"
	note+="<b><h2>Support Level:</b>" + "None"+"</h2>\n"
	// note+="<b>:</b>"
	//******************constructing integration string******************
	
	if integrationName != "" {
		note += "<b>Integration: </b>" + integrationName
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
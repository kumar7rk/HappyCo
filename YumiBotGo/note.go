package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
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
	note := makeNote(ID)

	addNote(conversationID, note)
}

//********************************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(ID string) (businessRec []Business, iapRec []IAP, integrationName string, planTypeRec []Plan) {
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
func makeNote(us_id string) (string) {
	var note string
	var formattedDate string

	//getting user data from the database
	businessRec, iapRec, integrationName, planTypeRec := getUserData(us_id)

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
		var BusinessPermission = business.PermissionsModel

		if BusinessPermission == "constant-roles" {
			BusinessPermission = "Contant/Full"
		}

		if BusinessPermission == "basic-roles" {
			BusinessPermission = "Basic"
		}
		note += "<b>Business: </b>" + business.Name + "\n"
		note += "<b>BusinessID:</b>" + business.ID + "\n"
		note += "<b>Permissions:</b>" + BusinessPermission + "\n"
		note += "<b>Role:</b>" + roles[permission] + "\n"
	}
	//******************constructing plan type string******************

	for _, plan := range planTypeRec {
		plan.Type = strings.Replace(plan.Type, "_", " ", -1)
		note += "<b>Plan: </b>" + strings.Title(plan.Type) + "\n"
	}
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
	return note
}

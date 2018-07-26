/* Adding User Data from DB to Intercom as a note for every new conversation

 */

package main

//adding import statements
import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	intercom "gopkg.in/intercom/intercom-go.v2"
)

//********************************************Variable declaration********************************************

type Inspection struct {
	Business     string
	User         string
	Role         string
	FolderID     string `db:"folder_id"`
	FolderName   string `db:"folder_name"`
	CreatedAt    string `db:"created_at"`
	TemplateName string `db:"template_name"`
	ID           string
	Status       string
	Location     string
}

type Report struct {
	Business   string
	User       string
	Role       string
	FolderID   string `db:"folder_id"`
	FolderName string `db:"folder_name"`
	CreatedAt  string `db:"created_at"`
	Name       string
	PublicID   string `db:"public_id"`
	Location   string
}

type Business struct {
	ID   string `db:"business_id"`
	Role string `db:"business_role_id"`
}
type IAP struct {
	Expiry string `db:"expires_at"`
}

type Integration struct {
	ID int `db:"id"`
}

type Plan struct {
	Type string `db:"plan_type"`
}

var plan_type_replica string

// date and time formatted string
var formattedDate string

// integration's name
var integrationName string

//recs
var inspectionsRec []Inspection

var reportsRec []Report

var businessRec []Business

var iapRec []IAP

var integrationRec []Integration

var planRec []Plan

var db *sqlx.DB

// structs for reading payload in json received from Intercom
type User struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

type Item struct {
	ConversationID string `json:"id"`
	User           User   `json:"user"`
}

type Data struct {
	Item Item `json:"item"`
}

type Message struct {
	Data Data `json:"data"`
}

//********************************************Main function********************************************

func main() {
	// loading .env file
	err := loadEnv()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	postgresURI := formURI()
	if postgresURI == "" {
		fmt.Println("URI error")
	}

	db, err = connect(postgresURI)

	if err == nil {
		fmt.Println("DB Connected")
	}
	if err != nil {
		fmt.Println("Error connecting")
	}

	//handling every new convesations in newConversation method
	http.HandleFunc("/", newConversation)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

//********************************************New Conversation********************************************

//gets intercom token, admin list, reads the payload, and post note as a reply in the conversation
func newConversation(w http.ResponseWriter, r *http.Request) {

	// Read body/payload
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unmarshal the json
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* getting attributes from the received json
	user type - lead/user
	userId - happyCo user id
	conversationId - Intercom conversation Id
	*/

	userType := msg.Data.Item.User.Type
	userId := msg.Data.Item.User.UserID            //65135
	conversationId := msg.Data.Item.ConversationID //15363702969

	//only run the following code when the received message is from a HappyCo user
	if userType == "user" {	
		go makeAndSendNote(userId, conversationId)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}

func makeAndSendNote(ID string, conversationID string) {
	p := fmt.Println

	// gets intercom access token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN")
	ic := intercom.NewClient(accessToken, "")

	user, err := ic.Users.FindByUserID(ID)
	if err != nil {
	   fmt.Fprintf(os.Stderr, "Error while finding user ID %v: %v\n", ID, err)
	   return
 	}	
 
 	//testing prints
	p("Conversation id: " + conversationID)
	p("User id: " + ID)
	p("User name: " + user.Name)

	// getting admin list from Intercom
	adminList, err := ic.Admins.List()
	if err != nil {
	   fmt.Fprintf(os.Stderr, "Error while loading admin list %v: %v\n", "", err)
	   return
 	}
	admins := adminList.Admins

	// setting admin to HappyBot
	// Adds the note from admin HappyBot
	admin := admins[13] 

	// calling the method to compile the note with all the required information
	note := makeNote(ID)

	ic.Conversations.Reply(conversationID, &admin, intercom.CONVERSATION_NOTE, note)
	//copied and pasted from api-docs
	if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
		fmt.Fprintf(os.Stderr, "Error from Intercom when replying %v: %v\n", "", err)
	   	return
	}
	var buildiumMessage string
	//extracting firstName from the user name.
	firstName := strings.Fields(user.Name)

	buildiumMessage = "Hi " + firstName[0] + "  \n \n Buildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫  \n <b>Our friends at Buildium support your Happy Inspector subscription and mobile app and can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>   \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyCo team ☺"

	if plan_type_replica == "buildium" {
		//ic.Conversations.Reply(conversationID,&admin,intercom.CONVERSATION_COMMENT,"Testing on internal plan \n " +buildiumMessage)
		plan_type_replica = buildiumMessage //fun stuff
		plan_type_replica = ""
	}
	
}

//********************************************Loading Enviornment Variables********************************************

// loading env file to load db parameters
func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

//********************************************Forming DB URI********************************************

// forming postgres URI
// returns string
func formURI() (str string) {

	// defining db parameters
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// buildiing db connection string
	postgresURI := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	// returing uri
	return postgresURI
}

//********************************************Connecting with DB********************************************

//connecting to the database
//returne db- sqlx
func connect(postgresURI string) (*sqlx.DB, error) {
	//opening connection using sqlx package
	db, err := sqlx.Connect("postgres", postgresURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//********************************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(u_id string) {

	//fetching most recent (5) inspections for the user within the last 30 days.

	err := db.Select(&inspectionsRec, "SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at as created_at,i.template_name,i.id,i.status,i.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $1::varchar AND i.archived_at IS NULL AND i.created_at > (CURRENT_DATE- interval '30 day') ORDER BY i.created_at DESC LIMIT $2", u_id, 5)
	if err != nil {
	   fmt.Fprintf(os.Stderr, "Error in inspection query %v: %v\n", u_id, err)
	   	return
	}

	//fetching most recent (5) reports for the user within the last 30 days.

	err = db.Select(&reportsRec, "SELECT folders.business,folders.user,folders.role,folders.folder_id,folders.folder_name,r.created_at as created_at,r.name,r.public_id,r.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN reports_v3 as r ON folders.folder_id = r.folder_id WHERE r.user_id = $1::varchar AND r.archived_at IS NULL AND r.created_at > (CURRENT_DATE- interval '30 day') ORDER BY r.created_at DESC LIMIT $2", u_id, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reoprt query %v: %v\n", u_id, err)
	   	return
	}
	// fetching business id and role id for user role in this business

	err = db.Select(&businessRec, "SELECT business_id,business_role_id FROM business_membership WHERE user_id = $1 AND inactivated_at IS NULL", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in business query %v: %v\n", u_id, err)
	   	return
	}
	// checking if the business is on IAP get the expiry date
	err = db.Select(&iapRec, "SELECT expires_at FROM iap_receipts WHERE company_id IN (SELECT business_id FROM business_membership WHERE user_id = $1) ORDER BY expires_at DESC limit 1", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in inspection query %v: %v\n", u_id, err)
	   	return
	}

	// Check if the business has integration w/Yardi
	err = db.Select(&integrationRec, "Select id FROM integration_yardi_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Yardi %v: %v\n", u_id, err)
	   	return
	}
	for _, integration := range integrationRec {
		if integration.ID > 0 {
			integrationName = "Yardi"
		}
	}
	// MRI
	err = db.Select(&integrationRec, "Select id FROM integration_mri_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-MRI %v: %v\n", u_id, err)
	   	return
	}
	for _, integration := range integrationRec {
		if integration.ID > 0 {
			integrationName = "MRI"

		}
	}
	// Resman
	err = db.Select(&integrationRec, "Select id FROM integration_resman_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in integration query-Resman %v: %v\n", u_id, err)
	   	return
	}
	for _, integration := range integrationRec {
		if integration.ID > 0 {
			integrationName = "Resman"
		}
	}

	// DD/buildium/mri
	err = db.Select(&planRec, "Select plan_type FROM subscriptions WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)", u_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in plan query %v: %v\n", u_id, err)
	   	return
	}
	defer db.Close()

}

//********************************************Building note********************************************

// code starts running from here.
// build the note in a string format
// should be called when a new intercom message is received
func makeNote(us_id string) string {
	var note string

	// p:= fmt.Println
	//getting user data from the database
	getUserData(us_id)

	//******************constructing business string******************
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"

	note += "<b>✅User is associated with the following businesses</b><br/><br/>"

	// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
	// 8, 9 = Basic Admin, PM
	var roles = [10]string{"", "Admin", "Process Manger", "Inspector", "Limited Inspector", "", "", "", "Admin", "Process Manager"}

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
	integrationName = ""

	//******************constructing plan type string******************

	for _, plan := range planRec {

		if plan.Type == "due_diligence" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Due Diligence"
		}
		if plan.Type == "buildium" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "Buildium"
		}
		if plan.Type == "mri" {
			note += "\n"
			note += "\n"
			note += "Plan: " + "MRI"
		}
		plan_type_replica = plan.Type
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
		formattedDate = ""
		url = ""
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
		formattedDate = ""
		url = ""
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
	formattedDate = ""

	return note
}

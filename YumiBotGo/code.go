/* Code runs with ngrok (for testing) when a new conversation is received. :)

need to setup ngrok beforehand and point the webhook to the right ngrok url manually everytime your machine boot up

*/

package main

//adding import statements
import (
	"fmt"
	"strings"
	"os"
    "strconv"

   "time"
  
  "net/http"  
  "encoding/json"
  "io/ioutil"
	
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "github.com/joho/godotenv"
  
  intercom "gopkg.in/intercom/intercom-go.v2"

)

//********************************************Variable declaration********************************************
	// testing for moving to sqlx structScan
	type data struct{
		business_I string 
	  	user_id_I string
	  	role_I string
	  	folder_id_I string
	  	folder_name_I string
	  	created_at_I string
	  	template_name_I string
		id_I string
	  	status_I string
	  	location_I string
	}

	//main attributes
	type row struct{
		inspection string
		report string
		role string
		iap string
	}

	//db variables for inspections query
  	var business_I string
  	var user_id_I string
  	var role_I string
  	var folder_id_I string
  	var folder_name_I string
  	var created_at_I string
  	var template_name_I string
	var id_I string
  	var status_I string
  	var location_I string

  	//db variables for report query
  	var business_R string
  	var user_id_R string
  	var role_R string
  	var folder_id_R string
  	var folder_name_R string
  	var created_at_R string
  	var name_R string
	var public_id_R string
  	var location_R string

  	//db variable for role query
  	var business_ro string
  	var role_ro string

  	//db variable for iap query
  	var expires_at_iap string
	
	//db variable for integration query
	var business_integration int

	//array for inspections
	var r [5] row
	//array for reports
	var r1 [5] row
	//array for role
	var r2 [1] row
	//array for iap
	var r3 [1] row

	//var d [10] data

	//counter to add values in array r for inspections
	var inspectionCounter int
	//counter to add values in array r for reports
	var reportCounter int
	//counter to add values in array r for role
	var roleCounter int
	
	// note string to be displayed in intercom
	var note string 

	// date and time formatted string
	var formattedDate string

	// name of the integration if there's one
	var integration string
	

// structs for reading payload in json received from Intercom
// Could have better names; for sure :)
type User struct {
	UserID string `json:"user_id"`
	Type string `json:"type"`
}

type Item struct {
	ConversationID string `json:"id"`
	User User `json:"user"`
}

type Data struct {
	Item Item `json:"item"`
}

type Message struct {
	Data Data `json:"data"`
}

//********************************************Main function********************************************

//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
func main() {
	// loading .env file
	err:=loadEnv();
	if err!=nil {
		fmt.Println("Error loading .env file")
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
	p:= fmt.Println
 	w.Write([]byte("Received"))
 	// gets intercom access token
	accessToken := os.Getenv("INTERCOM_ACCESS_TOKEN") // change INTERCOM_ACCESS_TOKEN_TEST
	ic := intercom.NewClient(accessToken, "")


	// Read body/payload
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal the json
    var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/* getting attributes from the received json
	user type - lead/user
	userId - happyCo user id
	conversationId - Intercom conversation Id
	*/
	
	userType := msg.Data.Item.User.Type
	userId := msg.Data.Item.User.UserID //65135
	conversationId := msg.Data.Item.ConversationID //15363702969


	//only run the following code when the received message is from a HappyCo user
	if userType == "user" {
		user, err := ic.Users.FindByUserID(userId)
		_=err
		//testing prints
		p("Conversation id: "+ conversationId)
		p("User id: "+ userId)
		p("User name: "+user.Name)

		// getting admin list from Intercom
		adminList, err1 := ic.Admins.List()
		_=err1
		admins := adminList.Admins

		// setting admin to HappyBot
		// Adds the note from admin - HappyBot
		admin:=admins[13] // change [0]

		// calling the method to compile the note with all the required information
		noteBuilder(userId)

		ic.Conversations.Reply(conversationId,&admin,intercom.CONVERSATION_NOTE,note)
		//fmt.Println(convo)
	}	
}

//********************************************Loading Enviornment Variables********************************************

// loading env file to load db parameters
func loadEnv() (error){
	err := godotenv.Load(".env")
	if err != nil {
	  return err;
	}	

	return nil;
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
    return postgresURI;
}
//********************************************Connecting with DB********************************************

//connecting to the database
//returne db- sqlx
func connect(postgresURI string) (*sqlx.DB, error){
	//opening connection using sqlx package
	db, err := sqlx.Connect("postgres", postgresURI)
	if err != nil {
		return nil, err;
	}
	return db,nil;
}

//********************************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(u_id string){
	inspectionCounter = 0
	roleCounter = 0
	reportCounter = 0;
	
	postgresURI := formURI();
	if postgresURI=="" {
		fmt.Println("URI error")
	}

	db,err := connect(postgresURI);

	if err == nil {
		fmt.Println("DB Connected");
	}
	if err != nil {
		fmt.Println("Error connecting");
	}
	// data1:=data{}
	//var data1 data
	//var d[10] data
	//running the query on the db
	rows, err := db.Queryx("SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at as created_at,i.template_name,i.id,i.status,i.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $3 AND i.archived_at IS NULL AND i.created_at > (CURRENT_DATE- interval '30 day') ORDER BY i.created_at DESC LIMIT $2",u_id,5,u_id);
	if err != nil {
    	panic(err)
    }
  // fetching all the records 
  for rows.Next() {
  	
  	err = rows.Scan(&business_I, &user_id_I, &role_I, &folder_id_I,&folder_name_I, &created_at_I, &template_name_I, &id_I,&status_I, &location_I)
  	// err = rows.StructScan(&data1)
    // fmt.Println(data1)
  	if err != nil {
      panic(err)
    }

    r[inspectionCounter].inspection = id_I+ " " + folder_id_I+ " " + created_at_I+ " " + template_name_I
    // r[count].inspection = data1.id_I+ " " + data1.folder_id_I+ " " + data1.created_at_I+ " " + data1.template_name_I
    inspectionCounter++
  }
  err = rows.Err()
  if err != nil {
  	panic(err) 
  }

  rows1, err1 := db.Queryx("SELECT folders.business,folders.user,folders.role,folders.folder_id,folders.folder_name,r.created_at as created_at,r.name,r.public_id,r.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN reports_v3 as r ON folders.folder_id = r.folder_id WHERE r.user_id = $3 AND r.archived_at IS NULL AND r.created_at > (CURRENT_DATE- interval '30 day') ORDER BY r.created_at DESC LIMIT $2",u_id,5,u_id);
	if err1 != nil {
    	panic(err1)
    }
  // fetching all the records 
  for rows1.Next() {
  	
	err1 = rows1.Scan(&business_R, &user_id_R, &role_R, &folder_id_R,&folder_name_R, &created_at_R, &name_R, &public_id_R,&location_R)
  	//err1 = rows1.StructScan(&data1)
    //fmt.Println(data1)
  	if err1 != nil {
      panic(err1)
    }
     r[reportCounter].report = public_id_R+ " " + created_at_R+ " " + name_R
    //r[count].inspection = data1.id_I+ " " + data1.folder_id_I+ " " + data1.created_at_I+ " " + data1.template_name_I
    reportCounter++
  }
  err1 = rows1.Err()
  if err1 != nil {
  	panic(err1) 
  }

  rows2, err2 := db.Queryx("SELECT business_id,business_role_id FROM business_membership WHERE user_id = $1 AND inactivated_at IS NULL",u_id);
  if err2 != nil {
    	panic(err2)
  }
  for rows2.Next() {
	err2 = rows2.Scan(&business_ro, &role_ro)
	if err2 != nil {
    	panic(err2)
  	}
  	r[roleCounter].role = business_ro+ " " + role_ro
  	roleCounter++
  }
  err2 = rows2.Err()
  if err2 != nil {
  	panic(err2) 
  }


  rows3, err3 := db.Queryx("SELECT expires_at FROM iap_receipts WHERE company_id IN (SELECT business_id FROM business_membership WHERE user_id = $1) ORDER BY expires_at DESC limit 1",u_id);
  if err3 != nil {
    	panic(err3)
  }
  for rows3.Next() {
	err3 = rows3.Scan(&expires_at_iap)
	if err3 != nil {
    	panic(err3)
  	}
  	r[0].iap = expires_at_iap
  }
  err3 = rows3.Err()
  if err3 != nil {
  	panic(err3) 
  }

  //Yardi
  rows4, err4 := db.Queryx("Select count(id) FROM integration_yardi_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)",u_id);
  if err4 != nil {
    	panic(err4)
  }
  for rows4.Next() {
	err4 = rows4.Scan(&business_integration)
	if err4 != nil {
    	panic(err4)
  	}
  	if business_integration > 0 {
  		integration = "Yardi"
  		business_integration = 0;
  	}
  }
  err4 = rows4.Err()
  if err4 != nil {
  	panic(err4) 
  }

  // MRI
  rows5, err5 := db.Queryx("Select count(id) FROM integration_mri_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)",u_id);
  if err5 != nil {
    	panic(err5)
  }
  for rows5.Next() {
	err5 = rows5.Scan(&business_integration)
	if err5 != nil {
    	panic(err5)
  	}
  	if business_integration > 0 {
  		integration = "MRI"
  		business_integration = 0;
  	}
  }
  err5 = rows5.Err()
  if err5 != nil {
  	panic(err5) 
  }

  // Resman
  rows6, err6 := db.Queryx("Select count(id) FROM integration_resman_properties WHERE business_id IN (SELECT business_id from business_membership WHERE user_id = $1 AND inactivated_at IS NULL)",u_id);
  if err6 != nil {
    	panic(err6)
  }
  for rows6.Next() {
	err6 = rows6.Scan(&business_integration)
	if err6 != nil {
    	panic(err6)
  	}
  	if business_integration > 0 {
  		integration = "Resman"
  		business_integration = 0;
  	}
  }
  err6 = rows6.Err()
  if err6 != nil {
  	panic(err6) 
  }

  defer db.Close()

}

//********************************************Building note********************************************

// code starts running from here.
// build the note in a string format
// should be called when a new intercom message is received
func noteBuilder(us_id string) {
	// p:= fmt.Println
	//getting user data from the database
	getUserData(us_id);
	
//******************constructing business string******************
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"

	note += "<b>✅User is associated with the following businesses</b><br/><br/>"

// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
// 8, 9 = Basic Admin, PM
	var roles = [10]string{"","Admin","Process Manger", "Inspector","Limited Inspector", "","","","Admin","Process Manager"}
	for i := 0; i < roleCounter; i++ {	
		split := strings.Fields(r[i].role)
		permission, err := strconv.Atoi(split[1])

		if split[1]=="1"||split[1]=="2"||split[1]=="3"||split[1]=="4" {
			note+=" The business is on Constant/full Permissions \n"
		}

		if split[1]=="8"||split[1]=="9" {
			note+=" The business is on Basic Permissions \n"
		}

		if err!=nil {
			panic(err)
		}
		note += roles[permission] + " for "
		var text = "https://manage.happyco.com/admin/businesses/"+split[0]
		note+=text +"<br/><br/>"
	}

	note +="\n"
	note +="\n"

//******************constructing inspection string******************
  	note += "<b>✅   Yumi found these recent (max: 5) <em>Inspections in last 30 days:</em></b><br/>"
	note +="\n"

	for i := 0; i < inspectionCounter; i++ {
		split := strings.Fields(r[i].inspection)
		 var url = "https://manage.happyco.com/folder/"+split[1]+"/inspections/"+split[0]

		var date, _  = time.Parse(time.RFC3339, split[2])
		formats := []map[string]string{
				{"format": "02", "description": "Day"},
				{"format": "Jan", "description": "Month"},
				{"format": "2006", "description": "Year"},
				
				{"format": "3", "description": "Hours"},		
				{"format": "04", "description": "Minutes"},		
				{"format": "PM", "description": "AM or PM"}}

		for _, f := range formats {
			formattedDate += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formattedDate = strings.TrimSpace(formattedDate)
				formattedDate+=":"
			}
		}
		note += "<a href="+url+">"+url+"</a>" + " " + formattedDate
	    note +="\n"
		formattedDate =""
	}

//******************constructing report string******************

    note +="\n"
	note +="\n"
  	note += "<b>✅   Yumi found these recent (max: 5) <em>Reports in last 30 days:</em></b><br/>"
    note +="\n"
  	
	for i := 0; i < reportCounter; i++ {

		split := strings.Fields(r[i].report)
		var url = "https://manage.happyco.com/reports/"+split[0]
		var name string
		for i := 2; i < len(split); i++ {
			name += split[i]
			name += " "
		 }

		var date, _  = time.Parse(time.RFC3339, split[1])
		formats := []map[string]string{
				{"format": "02", "description": "Day"},
				{"format": "Jan", "description": "Month"},
				{"format": "2006", "description": "Year"},
				
				{"format": "3", "description": "Hours"},		
				{"format": "04", "description": "Minutes"},		
				{"format": "PM", "description": "AM or PM"}}

		for _, f := range formats {
			formattedDate += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formattedDate = strings.TrimSpace(formattedDate)
				formattedDate+=":"
			}
		}
		 note += "<a href="+url+">"+name+"</a>" + " " + formattedDate
		 note +="\n"
		 formattedDate=""
	}
	var date, _  = time.Parse(time.RFC3339, expires_at_iap)
		formats := []map[string]string{
				{"format": "02", "description": "Day"},
				{"format": "Jan", "description": "Month"},
				{"format": "2006", "description": "Year"},
				
				{"format": "3", "description": "Hours"},		
				{"format": "04", "description": "Minutes"},		
				{"format": "PM", "description": "AM or PM"}}

		for _, f := range formats {
			formattedDate += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formattedDate = strings.TrimSpace(formattedDate)
				formattedDate+=":"
			}
		}
//******************constructing iap string******************
// this is only a part of the code if the business is on iap
	if expires_at_iap != "" {
		note+="\n"
		note+="\n"
		note+= "The business is on IAP. It expires on "+formattedDate
	}
		 formattedDate=""

//******************constructing integration string******************
		 if integration != "" {
		 	note+="\n"
			note+="\n"
			note+= "The business is "+integration
		 }

		 //return note
}
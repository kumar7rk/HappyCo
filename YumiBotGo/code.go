// Code runs when a new conversation is received. :)
// need to setup ngrok beforehand and point the webhook to the right ngrok url manually

package main

//adding import statements
import (
	"fmt"
	"strings"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"

	"github.com/joho/godotenv"
    "os"
    "strconv"

   "time"

   	intercom "gopkg.in/intercom/intercom-go.v2"

  "net/http"
  
  "encoding/json"
 "io/ioutil"
)

//****************************Variable declaration********************************************
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
	var inspection_counter int
	//counter to add values in array r for reports
	var report_counter int
	//counter to add values in array r for role
	var role_counter int
	
	// note string to be displayed in intercom
	var note string 

	// date and time formatted string
	var formatted_date string

	// name of the integration if there's one
	var integration_is string
	

type Innermost struct {
	Key0 string `json:"user_id"`
	Key01 string `json:"type"`
}

type Inner struct {
	Key1 string `json:"id"`
	Key2 Innermost `json:"user"`
}

type Outer struct {
	Key3 Inner `json:"item"`
}

type Outmost struct {
	Key4 Outer `json:"data"`
}

//****************************Main function********************************************

//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
func main() {
	// loading .env file
	err:=loadEnv();
	if err!=nil {
		fmt.Println("Error loading .env file")
	}
	
	http.HandleFunc("/", newConversation)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
    return
  }
}

func newConversation(w http.ResponseWriter, r *http.Request) {
	p:= fmt.Println
 	w.Write([]byte("Received"))

	access_token := os.Getenv("INTERCOM_ACCESS_TOKEN") // change INTERCOM_ACCESS_TOKEN_TEST
	ic := intercom.NewClient(access_token, "")


	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
    var msg Outmost
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var user_type = msg.Key4.Key3.Key2.Key01
	user_id_In := msg.Key4.Key3.Key2.Key0 //65135
	var conversation_id_In = msg.Key4.Key3.Key1 //15363702969

	//var conversation_id_In = "15363702969"
	// user_id_In := "73252"

	
	if user_type == "user" {
		user, err := ic.Users.FindByUserID(user_id_In)
		_=err
		
	
		p("Conversation id: "+ conversation_id_In)
		p("User id: "+ user_id_In)
		p("User name: "+user.Name)


		adminList, err := ic.Admins.List()
		admins := adminList.Admins

		// setting admin to HappyBot
		// Adds the note from user named HappyBot
		
		admin:=admins[13] // change [0]

		noteBuilder(user_id_In)

		convo, err:= ic.Conversations.Reply(conversation_id_In,&admin,intercom.CONVERSATION_NOTE,note)
		//fmt.Println(convo)
		_=err
		_=convo
	}
}
//****************************Loading Enviornment Variables********************************************

// loading env file to load db parameters
func loadEnv() (error){
	err := godotenv.Load(".env")
	if err != nil {
	  return err;
	}	

	return nil;
}
//****************************Forming DB URI********************************************

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
//****************************Connecting with DB********************************************

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

//****************************Getting UserData********************************************

//queries the db and adds returned values in array
func getUserData(u_id string){
	inspection_counter = 0
	role_counter = 0
	report_counter = 0;
	
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

    r[inspection_counter].inspection = id_I+ " " + folder_id_I+ " " + created_at_I+ " " + template_name_I
    // r[count].inspection = data1.id_I+ " " + data1.folder_id_I+ " " + data1.created_at_I+ " " + data1.template_name_I
    inspection_counter++
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
     r[report_counter].report = public_id_R+ " " + created_at_R+ " " + name_R
    //r[count].inspection = data1.id_I+ " " + data1.folder_id_I+ " " + data1.created_at_I+ " " + data1.template_name_I
    report_counter++
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
  	r[role_counter].role = business_ro+ " " + role_ro
  	role_counter++
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
  		integration_is = "Yardi"
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
  		integration_is = "MRI"
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
  		integration_is = "Resman"
  		business_integration = 0;
  	}
  }
  err6 = rows6.Err()
  if err6 != nil {
  	panic(err6) 
  }

  defer db.Close()

}
//****************************Building note********************************************

// code starts running from here.
// build the note in a string format
// should be called when a new intercom message is received
func noteBuilder(us_id string) {
	// p:= fmt.Println
	//getting user data from the database
	getUserData(us_id);
	
//***************working to construct the business string***********
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"

	note += "<b>✅User is associated with the following businesses</b><br/><br/>"

// 1, 2, 3, 4 = Constant Admin, PM, Inspector, Limited Inspector
// 8, 9 = Basic Admin, PM
	var roles = [10]string{"","Admin","Process Manger", "Inspector","Limited Inspector", "","","","Admin","Process Manager"}
	for i := 0; i < role_counter; i++ {	
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

//***************working to construct the inspection string***********
  	note += "<b>✅   Yumi found these recent (max: 5) <em>Inspections in last 30 days:</em></b><br/>"
	note +="\n"

	for i := 0; i < inspection_counter; i++ {
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
			formatted_date += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formatted_date = strings.TrimSpace(formatted_date)
				formatted_date+=":"
			}
		}
		note += "<a href="+url+">"+url+"</a>" + " " + formatted_date
	    note +="\n"
		formatted_date =""
	}

//***************working to construct the report string***********

    note +="\n"
	note +="\n"
  	note += "<b>✅   Yumi found these recent (max: 5) <em>Reports in last 30 days:</em></b><br/>"
    note +="\n"
  	
	for i := 0; i < report_counter; i++ {

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
			formatted_date += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formatted_date = strings.TrimSpace(formatted_date)
				formatted_date+=":"
			}
		}
		 note += "<a href="+url+">"+name+"</a>" + " " + formatted_date
		 note +="\n"
		 formatted_date=""
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
			formatted_date += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				formatted_date = strings.TrimSpace(formatted_date)
				formatted_date+=":"
			}
		}
//***************working to construct the iap string***********
// this is only a part of the code if the business is on iap
	if expires_at_iap != "" {
		note+="\n"
		note+="\n"
		note+= "The business is on IAP. It expires on "+formatted_date
	}
		 formatted_date=""

//***************working to construct the integration string***********
		 if integration_is != "" {
		 	note+="\n"
			note+="\n"
			note+= "The business is "+integration_is
		 }
	//p(note)
}
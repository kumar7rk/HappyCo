
package main

//adding import statements
import (
	"fmt"
	"strings"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"

	"github.com/joho/godotenv"
    "os"

   "time"
)

type row struct{
	//main attributes
	role int
	inspection string
	report string
	iap int
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
  	
	// array for 10 rows just a number should be more in released version
	var r [10] row

	//counter to add values in array r
	var count int

	// note string to be displayed in intercom
	var note string 

//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
		var dataAndTime string
func main() {
	// loading .env file
	err:=loadEnv();
	if err!=nil {
		fmt.Println("Error loading .env file")
	}
	
	// calling noteBuilder
	//noteBuilder("65135");

	
	//var date, _ = time.Parse("Time:Z07:00T15:04:05 Date:2006-01-02 ", "Time:-03:30T19:18:35 Date:2119-10-29")

	//defaultFormat := "2006-01-02 15:04:05 PM -07:00 Jan Mon MST"

	formats := []map[string]string{
		{"format": "2", "description": "Day"},
		{"format": "Jan", "description": "Month"},
		{"format": "2006", "description": "Year"},
		

		{"format": "3", "description": "Hours"},		
		{"format": "04", "description": "Minutes"},		
		{"format": "PM", "description": "AM or PM"}}

		for _, f := range formats {
			//dataAndTime += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
			//	dataAndTime+=":"
			}
			//fmt.Print(date.Format(f["format"])+" ")
		}
//		fmt.Println(dataAndTime)

noteBuilder("65135")
	
}

// loading env file to load db parameters
func loadEnv() (error){
	err := godotenv.Load(".env")
	if err != nil {
	  return err;
	}	

	return nil;
}
// forming postgres URI
// returns string
func formURI() (str string) {

	// defining db parameters
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// buildiing db connection string
  	psqlInfo := fmt.Sprintf("host=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, user, password, dbname)

	// returing uri
    return psqlInfo;
}

//connecting to the database
//returne db- sqlx
func connect(psqlInfo string) (*sqlx.DB, error){
	//opening connection using sqlx package
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err;
	}
	return db,nil;
}
//queries the db and adds returned values in array
func getUserData(u_id string,){
	psqlInfo := formURI();
	if psqlInfo=="" {
		fmt.Println("URI error")
	}

	db,err := connect(psqlInfo);

	if err == nil {
		fmt.Println("Connected");
	}
	if err != nil {
		fmt.Println("Error connecting");
	}
	//running the query on the db
	// right now it's fetching last 5 inspection for user_id 63135
	rows, err := db.Query("SELECT folders.business,folders.user,folders.role ,folders.folder_id,folders.folder_name,i.created_at::date as created_at,i.template_name,i.id,i.status,i.location FROM (SELECT businesses.business_id as business,businesses.user_id as user,role_id as role,folder_id as folder_id,folder_name as folder_name FROM (SELECT bm.business_id as business_id,bm.user_id as user_id,bm.business_role_id as role_id,f.id as folder_id,f.name as folder_name FROM business_membership as bm JOIN portfolios as f ON bm.business_id = f.business_id WHERE bm.user_id = $1 AND bm.inactivated_at IS NULL AND f.inactivated_at IS NULL) as businesses GROUP BY businesses.business_id,businesses.role_id,businesses.user_id,folder_id,folder_name ORDER BY businesses.business_id ) as folders JOIN inspections as i ON folders.folder_id = i.folder_id WHERE i.user_id = $3 AND i.archived_at IS NULL ORDER BY i.created_at DESC LIMIT $2",u_id,5,u_id);
	if err != nil {
    	panic(err)
    }
  defer rows.Close()
  // fetching all the records 
  for rows.Next() {
  	
  	err = rows.Scan(&business_I, &user_id_I, &role_I, &folder_id_I,&folder_name_I, &created_at_I, &template_name_I, &id_I,&status_I, &location_I)
  	if err != nil {
      panic(err)
    }

    r[count].inspection = id_I+ " " + folder_id_I+ " " + created_at_I+ " " + template_name_I
    count++
  }
  err = rows.Err()
  if err != nil {
  	panic(err) 
  }
  defer db.Close()

}
//not used was used to print individual values fetched from db
// might be used when calling new query :)
func printValues() {
	//fmt.Println(id,folder_id,created_at,template_name)
	for i := 0; i < count; i++ {
		fmt.Println(r[i].inspection)
		
	}
}
//code starts running from here.
// build the note in a string format
// should be called when a new intercom message is received
func noteBuilder(us_id string) {
	p:= fmt.Println
	
	getUserData(us_id);
	
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"
  	note += "<b>✅   Yumi found these recent <em>Inspections:</em></b><br/>"

	for i := 0; i < count; i++ {
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
			dataAndTime += date.Format(f["format"]+ " ");
			if f["description"] == "Hours" {
				dataAndTime = strings.TrimSpace(dataAndTime)
				dataAndTime+=":"
			}
		}
		//dataAndTime+="\n"
	     note += "<a href="+url+">"+url+"</a>" + " " + dataAndTime
	     note +="\n"

	     dataAndTime =""
	}
		//p(dataAndTime)

	p(note)
}

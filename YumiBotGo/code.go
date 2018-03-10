
package main

//adding import statements
import (
	"fmt"
	"strings"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"

	"github.com/joho/godotenv"
    "os"
)
type row struct{
	//main attributes
	role int
	inspection string
	report string
	iap int
}
	//db variables
	var id string
  	var folder_id string
  	var created_at string
  	var template_name string
	
	// array for 10 rows just a number should be more in released version
	var r [10] row

	//counter to add values in array r
	var count int

	// note string to be displayed in intercom
	var note string 


//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
func main() {
	// loading .env file
	err:=loadEnv();
	if err!=nil {
		fmt.Println("Error loading .env file")
	}
	
	// calling noteBuilder
	noteBuilder("65135");
	
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
	rows, err := db.Query("SELECT id, folder_id, created_at, template_name FROM inspections WHERE user_id = $1 AND FOLDER_ID IN ( SELECT portfolio_id FROM portfolio_access_controls WHERE user_id = $3) AND archived_at IS NULL ORDER BY created_at DESC LIMIT $2", u_id,5,u_id)
	if err != nil {
    	panic(err)
    }
  defer rows.Close()
  // fetching all the records 
  for rows.Next() {
  	
  	err = rows.Scan(&id, &folder_id, &created_at, &template_name)
  	if err != nil {
      panic(err)
    }

    r[count].inspection = id+ " " + folder_id+ " " + created_at+ " " + template_name
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

	getUserData(us_id);
	note = "<b>A small note from Yumi 🐶</b><br/><br/>"
  	note += "<b>✅   Yumi found these recent <em>Inspections:</em></b><br/>"

	for i := 0; i < count; i++ {

		split := strings.Fields(r[i].inspection)
		var url = "https://manage.happyco.com/folder/"+split[1]+"/inspections/"+split[0]
	    var date = split[2]
	    note += "<a href="+url+">"+url+"</a>" + " " + date
	    note +="\n"
	}
	fmt.Println(note)
}

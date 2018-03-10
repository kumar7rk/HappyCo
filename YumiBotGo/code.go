
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
	role int
	inspection string
	report string
	iap int
}
	var id string
  	var folder_id string
  	var created_at string
  	var template_name string
	
	var r [10] row

	var count int

	var note string 


//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
func main() {
	
	err:=loadEnv();
	if err!=nil {
		fmt.Println("Error loading .env file")
	}
	
	noteBuilder("65135");
	
}
func loadEnv() (error){
	// loading env file to load db parameters
	err := godotenv.Load(".env")
	if err != nil {
	  return err;
	}	

	return nil;
}
func formURI() (codee string) {

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

func connect(psqlInfo string) (*sqlx.DB, error){
	//opening connection using sqlx package
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err;
	}
	return db,nil;
}
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
func printValues() {
	//fmt.Println(id,folder_id,created_at,template_name)
	for i := 0; i < count; i++ {
		fmt.Println(r[i].inspection)
		
	}
}
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
//
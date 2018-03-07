
package main

//adding import statements
import (
	"fmt"
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
//starting main function
// at this moment all of the code is in main function
// from db connect to displaying results on console
func main() {
	// loading env file to load db parameters
	err := godotenv.Load(".env")
	if err != nil {
	  panic(err)
	}	
	getUserData("65135");
		printValues();
}

func getUserData(u_id string){

	// defining db parameters
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// buildiing db connection string
  	psqlInfo := fmt.Sprintf("host=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, user, password, dbname)

	//opening connection using sqlx package
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
	  panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
  		panic(err)
	}

	fmt.Println("Connected");

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
    // printing the fetched values

  }
  err = rows.Err()
  if err != nil {
  	panic(err)
  }
	
}

func printValues() {
	//fmt.Println(id,folder_id,created_at,template_name)
	for i := 0; i < count; i++ {
		fmt.Println(r[i].inspection)
		
	}
}
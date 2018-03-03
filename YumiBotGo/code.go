package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
    "os"
)
func main() {
	err := godotenv.Load(".env")
	if err != nil {
	  panic(err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

  psqlInfo := fmt.Sprintf("host=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
	  panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
  		panic(err)
	}

	fmt.Println("Connected");
	rows, err := db.Query("SELECT id, folder_id, created_at, template_name FROM inspections WHERE user_id = $1 AND FOLDER_ID IN ( SELECT portfolio_id FROM portfolio_access_controls WHERE user_id = $3) AND archived_at IS NULL ORDER BY created_at DESC LIMIT $2", 65135,5,65135)
	if err != nil {
    	panic(err)
    }

  defer rows.Close()
  for rows.Next() {
  	var id string
  	var folder_id int
  	var created_at string
  	var template_name string
  	err = rows.Scan(&id, &folder_id, &created_at, &template_name)
  	if err != nil {
      panic(err)
    }
    fmt.Println(id,folder_id,created_at,template_name)
  }
  err = rows.Err()
  if err != nil {
  	panic(err)
  }
}
// package main

// import (
// 	"database/sql"

// 	"github.com/gin-gonic/gin"
// )

// type DBUser struct{
// 	db *sql.DB
// 	table string
// 	DBName string
// }

// func main(){
// 	r := gin.Default()

// }

package main

import (
    "fmt"
    "log"
	"database/sql"
)
type DBUser struct{
	db *sql.DB
	table string
	DBName string
}
func main() {
    // Create a new database
    dbCreated, err := CreateDB("mydb")
    if err != nil {
        log.Fatal(err)
    }
    if dbCreated {
        fmt.Println("Database created")
    }

    // Connect to the database
    db, ok := DbConn("mydb")
    if !ok {
        log.Fatal("Could not connect to database")
    }
	if db != nil {
		fmt.Println("Connected to database")
	}

    // Create a new DBUser
    user := &DBUser{
        db:    db,
        table: "users",
    }

    // Create a new table
    tableFields := map[string]string{
        "Organisation": "int",
        "Item":         "text",
        "Email":        "text",
    }
    user.CreateTable(tableFields)
    fmt.Println("Table created")

    // Insert some data
    data := map[string]interface{}{
        "Organisation": 1,
        "Item":         "item1",
        "Email":        "email1@example.com",
    }
    user.InsertData(data)
    fmt.Println("Data inserted")

    // Get the data
    result := user.GetData()
    fmt.Println("Data:", result)
}
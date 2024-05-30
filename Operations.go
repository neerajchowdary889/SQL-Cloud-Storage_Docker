package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Organisation  int    `json:"organisation"`
	Item  		  string `json:"item"`
	Email         string `json:"email"`

}

var db *sql.DB
var dbMutex sync.Mutex

func DbConn(DBName string) (*sql.DB, bool) {
	DBName = DBName + ".db"
	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return db, true
}

func CreateDB(DBName string) (bool, error) {
    DBName = DBName + ".db"

    if _, err := os.Stat(DBName); err == nil {
        log.Println("Database already exists")
        return false, nil
    } else if os.IsNotExist(err) {
        db, err := os.Create(DBName)
        if err != nil {
            log.Println(err)
            return false, err
        }
        defer db.Close()
        return true, nil
    } else {
        log.Println(err)
        return false, err
    }
}

func(user *DBUser) CreateTable(tableFields map[string]string) (bool, error){
    // db := user.db
    fieldDefs := make([]string, 0, len(tableFields))
    fieldDefs = append(fieldDefs, "id INTEGER PRIMARY KEY AUTOINCREMENT")
    for field, fieldType := range tableFields {
        fieldDefs = append(fieldDefs, fmt.Sprintf("%s %s", field, strings.ToUpper(fieldType)))
    }

    fmt.Println(fieldDefs)

    table := fmt.Sprintf(
        "CREATE TABLE %s (%s);",
        user.table,
        strings.Join(fieldDefs, ", "),
    )

    fmt.Println(table)

    statement, err := (user.db).Prepare(table)
    if err != nil {
        return false, err
    }
    statement.Exec()
    return true, nil
}

func (user *DBUser) InsertData(data map[string]interface{}) {
    dbMutex.Lock()
    defer dbMutex.Unlock()

    fields := make([]string, 0, len(data))
    placeholders := make([]string, 0, len(data))
    values := make([]interface{}, 0, len(data))
    for field, value := range data {
        fields = append(fields, field)
        placeholders = append(placeholders, "?")
        values = append(values, value)
    }

    insert := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s);",
        user.table,
        strings.Join(fields, ", "),
        strings.Join(placeholders, ", "),
    )

    statement, err := user.db.Prepare(insert)
    if err != nil {
        log.Fatal(err.Error())
    }
    _, err = statement.Exec(values...)
    if err != nil {
        log.Fatal(err.Error())
    }
    log.Println("Data inserted")
}

func (user *DBUser) GetData() []map[string]interface{} {
    rows, err := (user.db).Query("SELECT * FROM " + user.table)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer rows.Close()

    var result []map[string]interface{}
    cols, _ := rows.Columns()

    for rows.Next() {
        columns := make([]interface{}, len(cols))
        columnPointers := make([]interface{}, len(cols))
        for i := range columns {
            columnPointers[i] = &columns[i]
        }

        if err := rows.Scan(columnPointers...); err != nil {
            log.Fatal(err.Error())
        }

        rowData := make(map[string]interface{})
        for i, colName := range cols {
            val := columnPointers[i].(*interface{})
            rowData[colName] = *val
        }

        result = append(result, rowData)
    }

    return result
}


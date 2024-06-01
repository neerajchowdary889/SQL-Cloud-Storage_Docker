package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var dbMutex sync.RWMutex

func DbConn(DBName string) (*sql.DB, bool) {
    dir := "MyDBs"
    DBName = filepath.Join(dir, DBName+".db")
    db, err := sql.Open("sqlite3", DBName)
    if err != nil {
        log.Println(err)
        return nil, false
    }

    return db, true
}

func CreateDB(DBName string) (bool, error) {
    dbMutex.Lock()
    defer dbMutex.Unlock()

    dir := "MyDBs"
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.Mkdir(dir, 0755)
    }

    DBName = filepath.Join(dir, DBName+".db")

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
    dbMutex.Lock()
    defer dbMutex.Unlock()

    fieldDefs := make([]string, 0, len(tableFields))
    fieldDefs = append(fieldDefs, "id INTEGER PRIMARY KEY AUTOINCREMENT")
    for field, fieldType := range tableFields {
        fieldDefs = append(fieldDefs, fmt.Sprintf("%s %s", field, strings.ToUpper(fieldType)))
    }

    table := fmt.Sprintf(
        "CREATE TABLE %s (%s);",
        user.table,
        strings.Join(fieldDefs, ", "),
    )

    statement, err := (user.db).Prepare(table)
    if err != nil {
        return false, err
    }
    statement.Exec()
    return true, nil
}

func (user *DBUser) InsertData(data map[string]interface{}) (bool, error) {
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
        return false, err
    }
    _, err = statement.Exec(values...)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (user *DBUser) GetData() []map[string]interface{} {
    dbMutex.RLock()
    defer dbMutex.RUnlock()

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

func (user *DBUser) UpdateData(data map[string]interface{}, condition map[string]interface{}) (bool, error){
    dbMutex.Lock()
    defer dbMutex.Unlock()

    fields := make([]string, 0, len(data))
    values := make([]interface{}, 0, len(data))
    for field, value := range data {
        fields = append(fields, fmt.Sprintf("%s = ?", field))
        values = append(values, value)
    }

    conditions := make([]string, 0, len(condition))
    for field, value := range condition {
        conditions = append(conditions, fmt.Sprintf("%s = ?", field))
        values = append(values, value)
    }

    update := fmt.Sprintf(
        "UPDATE %s SET %s WHERE %s;",
        user.table,
        strings.Join(fields, ", "),
        strings.Join(conditions, " AND "),
    )

    statement, err := (user.db).Prepare(update)
    if err != nil {
        return false, err
    }
    _, err = statement.Exec(values...)
    if err != nil {
        return false, err
    }
    log.Println("Data updated")
    return true, nil
}

// func (user *DBUser) ReadData(columns []string, condition map[string]interface{}) (*sql.Rows, error) {
//     dbMutex.RLock()
//     defer dbMutex.RUnlock()

//     conditions := make([]string, 0, len(condition))
//     values := make([]interface{}, 0, len(condition))
//     for field, value := range condition {
//         conditions = append(conditions, fmt.Sprintf("%s = ?", field))
//         values = append(values, value)
//     }

//     query := fmt.Sprintf(
//         "SELECT %s FROM %s WHERE %s;",
//         strings.Join(columns, ", "),
//         user.table,
//         strings.Join(conditions, " AND "),
//     )

//     statement, err := (user.db).Prepare(query)
//     if err != nil {
//         log.Fatal(err.Error())
//     }
//     rows, err := statement.Query(values...)
//     if err != nil {
//         log.Fatal(err.Error())
//     }

//     return rows, nil
// }

func (user *DBUser) ReadDataWithConditions(columns []string, condition []string) ([]map[string]interface{}, error) {
    dbMutex.RLock()
    defer dbMutex.RUnlock()

    combinedCondition := strings.Join(condition, " ")

    query := fmt.Sprintf(
        "SELECT %s FROM %s WHERE %s;",
        strings.Join(columns, ", "),
        user.table,
        combinedCondition,
    )
    fmt.Print(query)
    statement, err := user.db.Prepare(query)
    if err != nil {
        return nil, err
    }
    defer statement.Close()

    rows, err := statement.Query()
    if err != nil {
        return nil, err
    }

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
    return result, nil
}

func (user *DBUser) GetRangeData(limit int, offset int) []map[string]interface{} {
    dbMutex.RLock()
    defer dbMutex.RUnlock()

    query := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", user.table, limit, offset)
    rows, err := (user.db).Query(query)
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

func (user *DBUser) DeleteData(condition map[string]interface{}) (bool, error){
    dbMutex.Lock()
    defer dbMutex.Unlock()

    conditions := make([]string, 0, len(condition))
    values := make([]interface{}, 0, len(condition))
    for field, value := range condition {
        conditions = append(conditions, fmt.Sprintf("%s = ?", field))
        values = append(values, value)
    }

    delete := fmt.Sprintf(
        "DELETE FROM %s WHERE %s;",
        user.table,
        strings.Join(conditions, " AND "),
    )

    statement, err := (user.db).Prepare(delete)
    if err != nil {
        return false, err
    }
    _, err = statement.Exec(values...)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (user *DBUser) DropTable() (bool, error){
    dbMutex.Lock()
    defer dbMutex.Unlock()

    drop := fmt.Sprintf("DROP TABLE %s;", user.table)
    statement, err := (user.db).Prepare(drop)
    if err != nil {
        return false, err
    }
    _, err = statement.Exec()
    if err != nil {
        return false, err
    }
    return true, nil
}

func (user *DBUser) DropDB() (bool, error){
    dbMutex.Lock()
    defer dbMutex.Unlock()

    drop := fmt.Sprintf("DROP DATABASE %s;", user.DBName)
    statement, err := (user.db).Prepare(drop)
    if err != nil {
        return false, err
    }
    _, err = statement.Exec()
    if err != nil {
        return false, err
    }
    return true, nil
}

func (user *DBUser) AlterTable(addFields map[string]string, dropFields []string) (bool, error){
    dbMutex.Lock()
    defer dbMutex.Unlock()

    addFieldDefs := make([]string, 0, len(addFields))
    if len(addFields) != 0{
        for field, fieldType := range addFields {
            addFieldDefs = append(addFieldDefs, fmt.Sprintf("ADD %s %s", field, strings.ToUpper(fieldType)))
        }
    }

    dropFieldDefs := make([]string, 0, len(dropFields))
    if len(dropFields) != 0 {
        for _, field := range dropFields {
            dropFieldDefs = append(dropFieldDefs, fmt.Sprintf("DROP COLUMN %s", field))
        }
    }

    alter := fmt.Sprintf(
        "ALTER TABLE %s %s %s;",
        user.table,
        strings.Join(addFieldDefs, ", "),
        strings.Join(dropFieldDefs, ", "),
    )

    statement, err := (user.db).Prepare(alter)
    if err != nil {
        return false, err
    }
    _, err = statement.Exec()
    if err != nil {
        return false, err
    }
    return true, nil
}
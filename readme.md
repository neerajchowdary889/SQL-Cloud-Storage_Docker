# SQL Cloud Storage Docker

This project is a RESTful API server written in Go. It provides endpoints for managing a SQL database. The server uses the Gin framework for handling HTTP requests and Gorilla sessions for session management.

## Endpoints

### POST /init/CreateDB

Creates a new database. The request body should contain a JSON object with a `DBName` key.

### POST /init/ConnectDB

Connects to a database. The request body should contain a JSON object with a `DBName` key.

### POST /init/CreateTable

Creates a new table in the connected database. The request body should contain a JSON object with `tableName` and `tableFields` keys.

### POST /init/InsertData

Inserts data into the connected database. The request body should contain a JSON object with a `data` key.

### POST /init/UpdateData

Updates data in the connected database. The request body should contain a JSON object with `updateData` and `condition` keys.

### POST /init/DeleteData

Deletes data from the connected database. The request body should contain a JSON object with a `condition` key.

### GET /init/GetData

Retrieves all data from the connected database.

### GET /init/ReadDataWithCondition

Retrieves data from the connected database based on a condition. The request body should contain a JSON object with `columns` and `condition` keys.

### GET /init/GetRangeData

Retrieves a range of data from the connected database. The request body should contain a JSON object with `limit` and `offset` keys.

### POST /init/AlterTable

Alters a table in the connected database. The request body should contain a JSON object with `tableFields` and `dropFields` keys.

### POST /init/DropTable

Drops a table from the connected database.

### POST /init/DropDB

Drops the connected database.

## Running the Server

The server runs on port 5150. To start the server, run the `main.go` file.

**go** **run** **main.go**

## Dependencies

* Gin
* Gorilla sessions
* SQL database driver (depends on your SQL database)

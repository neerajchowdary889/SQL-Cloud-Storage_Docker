package main

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)
type DBUser struct{
	db *sql.DB
	table string
	DBName string
}


var store = sessions.NewCookieStore([]byte("Sql-Cloud-Storage_Docker"))


func main() {
    r := gin.Default()

    r.POST("/init/CreateDB", func(c *gin.Context){ 

        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        dbCreated, err := CreateDB(input["DBName"].(string))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if dbCreated {
            c.JSON(http.StatusOK, gin.H{"message": "Database created", "DBName": input["DBName"]})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not created"})
        return
    });
    
    r.POST("/init/ConnectDB", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        db, connected := DbConn(input["DBName"].(string))
        if !connected {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not connected"})
            return
        }
        session.Values["User"] = &DBUser{
            db:    db,
        }
        session.Save(c.Request, c.Writer)
        c.JSON(http.StatusOK, gin.H{"message": "Database connected", "DBName": input["DBName"]})
        return
    });
    
    r.POST("/init/CreateTable", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            
        }

        user := session.Values["User"].(*DBUser)
        user.table = input["tableName"].(string)
        tableFields := input["tableFields"].(map[string]string)
        tableCreated, err := user.CreateTable(tableFields)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if tableCreated {
            c.JSON(http.StatusOK, gin.H{"message": "Table created"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Table not created"})
        return
    });
    
    r.POST("/init/InsertData", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        user := session.Values["User"].(*DBUser)
        data := input["data"].(map[string]interface{})
        dataInserted, err := user.InsertData(data)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if dataInserted {
            c.JSON(http.StatusOK, gin.H{"message": "Data inserted"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not inserted"})
        return
    });
    r.POST("/init/UpdateData", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        user := session.Values["User"].(*DBUser)
        updateData := input["updateData"].(map[string]interface{})
        condition := input["condition"].(map[string]interface{})
        dataUpdated, err := user.UpdateData(updateData, condition)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if dataUpdated {
            c.JSON(http.StatusOK, gin.H{"message": "Data updated"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not updated"})
        return
    });

    r.POST("/init/DeleteData", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        user := session.Values["User"].(*DBUser)
        condition := input["condition"].(map[string]interface{})
        dataDeleted, err := user.DeleteData(condition)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if dataDeleted {
            c.JSON(http.StatusOK, gin.H{"message": "Data deleted"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Data not deleted"})
        return
    });

    r.GET("/init/GetData", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        user := session.Values["User"].(*DBUser)
        result := user.GetData()

        c.JSON(http.StatusOK, gin.H{"data": result})
        return
    });

    r.GET("/init/ReadDataWithCondition", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        user := session.Values["User"].(*DBUser)
        columns := input["columns"].([]string)
        condition := input["condition"].([]string)
        result, err := user.ReadDataWithConditions(columns, condition)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"data": result})
        return
    });

    r.GET("/init/GetRangeData", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user := session.Values["User"].(*DBUser)
        limit := input["limit"].(int)
        offset := input["offset"].(int)
        result := user.GetRangeData(limit, offset)
        c.JSON(http.StatusOK, gin.H{"data": result})
        return
    });

    r.POST("/init/AlterTable", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        var input map[string]interface{}

        if err := c.ShouldBindJSON(&input); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
        }

        user := session.Values["User"].(*DBUser)

        tableFields, ok := input["tableFields"].(map[string]string)
        if !ok || len(tableFields) == 0 {
            tableFields = nil
        }
        
        dropFields, ok := input["dropFields"].([]string)
        if !ok || len(dropFields) == 0 {
            dropFields = nil
        }

        tableAltered, err := user.AlterTable(tableFields, dropFields)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if tableAltered {
            c.JSON(http.StatusOK, gin.H{"message": "Table altered"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Table not altered"})
        return
    });

    r.POST("/init/DropTable", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        user := session.Values["User"].(*DBUser)
        tableDropped, err := user.DropTable()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if tableDropped {
            c.JSON(http.StatusOK, gin.H{"message": "Table dropped"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Table not dropped"})
        return
    });

    r.POST("/init/DropDB", func(c *gin.Context){
        session, err := store.Get(c.Request, "DB-details")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        user := session.Values["User"].(*DBUser)
        dbDropped, err := user.DropDB()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if dbDropped {
            c.JSON(http.StatusOK, gin.H{"message": "Database dropped"})
            return
        }
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not dropped"})
         return
    });

    r.Run(":5150")
}
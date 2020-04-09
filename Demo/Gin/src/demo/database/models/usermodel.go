package models

import (
	"database/sql"
	"demo/database/items"
	"fmt"

	"github.com/gin-gonic/gin"
)

type usermodel struct {
	DB *sql.DB
}

const (
	host     = "localhost"
	user     = "postgres"
	password = "123"
	dbname   = "mydb"
)

func connectDB() (db *sql.DB, err error) {

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err = sql.Open("postgres", conn)
	if err != nil {
		fmt.Printf("Fail to openDB: %v \n", err)
		return nil, err

	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Fail to conenct: %v \n", err)
		return nil, err
	}

	fmt.Println("Ping OK")
	return db, nil
}

func connectPostgreSQL() {

	db, err := connectDB()
	if err != nil {
		return
	}

	_sql := "SELECT user_name, full_name FROM public.account LIMIT 1;"
	row, err := db.Query(_sql)
	if err != nil {
		fmt.Printf("Fail to query: %v \n", err)
		return
	}
	var col1 string
	var col2 string
	for row.Next() {
		row.Scan(&col1, &col2)
		fmt.Printf("value Col1: %v \n", col1)
		fmt.Printf("value Col2: %v \n", col2)
	}
	fmt.Println("End !!!")
}

func GetInfoUser(c *gin.Context) {

	//Connect DB
	db, err := connectDB()
	if err != nil {
		return
	}
	//Get info
	sql := "SELECT user_name, full_name FROM public.account LIMIT 1;"
	row, err := db.Query(sql)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"messages": "Fail",
		})
		return
	}

	fmt.Println(sql)

	_infoUser := items.User{}
	var userName string
	var fullName string
	for row.Next() {
		row.Scan(&userName, &fullName)
	}
	_infoUser.UserName = userName
	_infoUser.FullName = fullName

	defer db.Close()

	fmt.Println(_infoUser)
	c.JSON(200, _infoUser)
}

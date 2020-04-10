package models

import (
	"database/sql"
	"demo/common"
	"demo/database/items"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userPostForm struct {
	userName string
	passWord string
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

	//Get post form
	infoUser := userPostForm{
		userName: fmt.Sprintf("'%v'", c.PostForm("user_name")),
		passWord: fmt.Sprintf("'%v'", c.PostForm("pass_word")),
	}

	fmt.Println(infoUser)
	//Get info
	sql := `SELECT user_name, full_name FROM public.user WHERE user_name = ` + infoUser.userName + ` AND pass_word = ` + infoUser.passWord
	row, err := db.Query(sql)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"messages": "Fail",
		})
		return
	}

	_infoUser := items.User{}
	var userName string
	var fullName string
	for row.Next() {
		row.Scan(&userName, &fullName)
	}
	if userName == "" {
		c.JSON(501, gin.H{
			"messages": common.MsgLoginError,
		})
		return
	}
	_infoUser.UserName = userName
	_infoUser.FullName = fullName

	defer db.Close()

	fmt.Println(_infoUser)
	c.JSON(200, _infoUser)
}

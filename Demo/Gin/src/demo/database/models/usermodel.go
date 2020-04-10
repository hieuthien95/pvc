package models

import (
	"database/sql"
	"demo/common"
	"demo/database/items"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userPostForm struct {
	userID   string
	userName string
	passWord string
	email    string
}

const (
	host     = "localhost"
	user     = "congpv"
	password = "Marumori1103"
	dbname   = "my_db"
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

// GetInfoUser Get info user with user_id
func GetInfoUser(c *gin.Context) {

	//Connect DB
	db, err := connectDB()
	if err != nil {
		return
	}

	//Get post form
	infoUser := userPostForm{
		userName: fmt.Sprintf("'%v'", c.PostForm("user_id")),
		passWord: fmt.Sprintf("'%v'", c.PostForm("password")),
	}

	fmt.Println(infoUser)

	//Get info
	sql := `SELECT username, email FROM public.account WHERE user_id = ` + infoUser.userName + ` AND password = ` + infoUser.passWord
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
	var email string
	for row.Next() {
		row.Scan(&userName, &email)
	}
	if userName == "" {
		c.JSON(501, gin.H{
			"messages": common.MsgLoginError,
		})
		return
	}
	_infoUser.UserName = userName
	_infoUser.Email = email

	defer db.Close()

	fmt.Println(_infoUser)
	c.JSON(200, _infoUser)
}

//UpdateInfoUser update password/email user with user_id
func UpdateInfoUser(c *gin.Context) {

	//Conenct DB
	db, err := connectDB()
	if err != nil {
		return
	}

	//Get info form
	infoUser := userPostForm{
		userID:   fmt.Sprintf("'%v'", c.PostForm("user_id")),
		passWord: fmt.Sprintf("'%v'", c.PostForm("password")),
		email:    fmt.Sprintf("'%v'", c.PostForm("email")),
	}

	fmt.Println(infoUser)

	//Create query udpate
	sql := ""

	if infoUser.email == "''" {
		sql = `UPDATE public.account
		SET password = ` + infoUser.passWord + ` , date_update=current_date
		WHERE user_id = ` + infoUser.userID
	} else if infoUser.passWord == "''" {
		sql = `UPDATE public.account
		SET email = ` + infoUser.email + ` , date_update=current_date
		WHERE user_id = ` + infoUser.userID
	} else {
		sql = `UPDATE public.account
		SET email = ` + infoUser.email + ` , password = ` + infoUser.passWord + `, date_update=current_date
		WHERE user_id = ` + infoUser.userID
	}

	if _, err = db.Exec(sql); err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"messages": "Update Fail",
		})
		return
	}

	c.JSON(200, gin.H{
		"messages": "Update Success",
	})

	defer db.Close()
}

// DeleteUser delete user with user_id
func DeleteUser(c *gin.Context) {

	//Conenct DB
	db, err := connectDB()
	if err != nil {
		return
	}

	//Get post form
	infoUser := userPostForm{
		userID: fmt.Sprintf("'%v'", c.PostForm("user_id")),
	}

	fmt.Println(infoUser)
	//Get info

	sql := `DELETE FROM public.account
	WHERE user_id = ` + infoUser.userID

	if _, err = db.Exec(sql); err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"messages": "Delete Fail",
		})
		return
	}

	c.JSON(200, gin.H{
		"messages": "Delete Success",
	})

	defer db.Close()
}

// InsertUser Insert info new user
func InsertUser(c *gin.Context) {

	//Conenct DB
	db, err := connectDB()
	if err != nil {
		return
	}

	//Get post form
	infoUser := userPostForm{
		userID:   fmt.Sprintf("'%v'", c.PostForm("user_id")),
		userName: fmt.Sprintf("'%v'", c.PostForm("username")),
		passWord: fmt.Sprintf("'%v'", c.PostForm("password")),
		email:    fmt.Sprintf("'%v'", c.PostForm("email")),
	}

	fmt.Println(infoUser)
	//Get info

	sql := `INSERT INTO public.account(
		user_id, username, password, email, date_create, date_update)
		VALUES (` + infoUser.userID + `, ` + infoUser.userName + `,` + infoUser.passWord + `,` + infoUser.email + `, current_date, current_date);`

	if _, err = db.Exec(sql); err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"messages": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"messages": "Insert Success",
	})

	defer db.Close()
}

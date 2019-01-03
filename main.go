package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"time"
)

type Note struct {
	Id    int
	Title string
	Text  string
	Date  time.Time
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "password"
	dbName := "notesdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	port := os.Getenv("PORT")
	//defer db.Close()

	//result, err := db.Exec("insert into notesdb.user_info (userName, description) values (?, ?)",
	//	"Hello", "Hello World!")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result.LastInsertId()) // id добавленного объекта
	//fmt.Println(result.RowsAffected()) // количество затронутых строк

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		log.Print("In GET")
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/note", func(c *gin.Context) {
		log.Print("In POST")
		db := dbConn()
		title := c.PostForm("title")
		description := c.PostForm("description")
		date := c.PostForm("date")
		log.Print("In POST title = " + title)
		log.Print("In POST description = " + description)
		log.Print("In POST date = " + date)
		insForm, err := db.Prepare("INSERT INTO note_info (title, text, date) VALUES(?, ?,?)")
		if err != nil {
			panic(err.Error())
			log.Print("In panic" + err.Error())
		}
		insForm.Exec(title, description, date)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
		//defer db.Close()
	})

	router.Run(":" + port)
}

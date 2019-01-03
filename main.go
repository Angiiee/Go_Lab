package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"time"
)

type Note struct {
	Id         int
	Title      string
	Text       string
	Date       time.Time
	DateString string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "password"
	dbName := "notesdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	port := os.Getenv("PORT")
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
		db := dbConn()
		rows, err := db.Query("SELECT * FROM notesdb.note_info ORDER BY note_info.date ASC")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		notes := []Note{}
		for rows.Next() {
			p := Note{}
			err := rows.Scan(&p.Id, &p.Title, &p.Text, &p.Date)
			if err != nil {
				fmt.Println(err)
				continue
			}
			p.DateString = p.Date.Format(time.RFC1123)
			notes = append(notes, p)
		}
		c.HTML(http.StatusOK, "index.tmpl.html", notes)
	})

	router.POST("/", func(c *gin.Context) {
		log.Print("In POST")
		db := dbConn()
		title := c.PostForm("title")
		description := c.PostForm("description")
		date := c.PostForm("date")
		insForm, err := db.Prepare("INSERT INTO note_info (title, text, date) VALUES(?, ?,?)")
		if err != nil {
			panic(err.Error())
			log.Print("In panic" + err.Error())
		}
		insForm.Exec(title, description, date)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
		defer db.Close()
	})

	router.POST("/delete", func(c *gin.Context) {
		log.Print("In DELETE")
		db := dbConn()
		id := c.PostForm("id")
		log.Print("In DELETE id =" + id)
		delForm, err := db.Prepare("DELETE FROM  notesdb.note_info WHERE id=?")
		if err != nil {
			panic(err.Error())
			log.Print("In panic" + err.Error())
		}
		delForm.Exec(id)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
		defer db.Close()
	})

	router.POST("/update", func(c *gin.Context) {
		log.Print("In UPDATE")
		db := dbConn()
		id := c.PostForm("idUpdate")
		title := c.PostForm("titleUpdate")
		description := c.PostForm("descriptionUpdate")
		date := c.PostForm("dateUpdate")
		updForm, err := db.Prepare("UPDATE note_info SET title=?, text=?, date=? WHERE id=?")
		if err != nil {
			panic(err.Error())
			log.Print("In panic" + err.Error())
		}
		updForm.Exec(title, description, date, id)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
		defer db.Close()
	})

	router.Run(":" + port)
}

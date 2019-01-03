package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	db, err := sql.Open("mysql", "root:password@/notesdb")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	//result, err := db.Exec("insert into notesdb.user_info (userName, description) values (?, ?)",
	//	"Hello", "Hello World!")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result.LastInsertId()) // id добавленного объекта
	//fmt.Println(result.RowsAffected()) // количество затронутых строк

	if port == "" {
		log.Fatal("$PORT must be set")
		//port = "8080"
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
		title := c.PostForm("title")
		description := c.PostForm("description")
		date := c.PostForm("date")
		log.Print("In POST title = " + title)
		log.Print("In POST description = " + description)
		log.Print("In POST date = " + date)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

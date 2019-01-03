package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/russross/blackfriday"
	"log"
	"net/http"
	"os"
)

var (
	repeat int
)

func repeatHandler(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!\n")
	}
	c.String(http.StatusOK, buffer.String())
}

func main() {
	port := os.Getenv("PORT")

	db, err := sql.Open("mysql", "root:password@/notesdb")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into notesdb.user_info (userName, description) values (?, ?)",
		"Hello", "Hello World!")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId()) // id добавленного объекта
	fmt.Println(result.RowsAffected()) // количество затронутых строк

	if port == "" {
		log.Fatal("$PORT must be set")
		//port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatHandler)

	router.Run(":" + port)
}

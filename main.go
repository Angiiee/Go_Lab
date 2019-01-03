package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	repeat int
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

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
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})

	router.Run(":" + port)
}

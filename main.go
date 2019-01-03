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
)

func processNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	log.Print("In processNote")
	if r.Method == "POST" {
		log.Print("In processNote POST")
		r.ParseForm()
		// logic part of log in
		title := r.FormValue("title")
		description := r.FormValue("description")
		date := r.FormValue("date")
		fmt.Fprintf(w, "title = %s\n", title)
		fmt.Fprintf(w, "description = %s\n", description)
		fmt.Fprintf(w, "date = %s\n", date)
		http.ServeFile(w, r, "index.tmpl.html")
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

	http.HandleFunc("/", processNote)

	router.GET("/", func(c *gin.Context) {
		log.Print("In GET")
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/", func(c *gin.Context) {
		log.Print("In POST")
		title := c.PostForm("title")
		log.Print("In POST title = " + title)
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

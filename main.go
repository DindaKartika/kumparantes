package main

import (
	"fmt"

	"kumparantes/databases"
	"kumparantes/models"
	"kumparantes/router"
)

func main() {
	fmt.Println("Welcome to the server")

	e := router.New()

	// init database
	autoCreateTables()
	autoMigrateTables()

	e.Start(":8000")
}

// autoCreateTables: create database tables using GORM
func autoCreateTables() {
	if !databases.App.DBConfig.HasTable(&models.Article{}) {
		databases.App.DBConfig.CreateTable(&models.Article{})
		if databases.App.ENV == "dev" {
			var articles []models.Article = []models.Article{
				models.Article{Author: "Dinda Kartika Vilaili", Title: "What is Lorem Ipsum?", Body: "Lorem Ipsum is simply dummy text of the printing and typesetting industry."},
				models.Article{Author: "Dinda Kartika Vilaili", Title: "Why do we use Lorem Ipsum?", Body: "The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English."},
				models.Article{Author: "Dinda Kartika Vilaili", Title: "Where can I get some?", Body: "There are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable."},
			}

			for _, article := range articles {
				databases.App.DBConfig.Create(&article)
			}
		}
	}
}

// autoMigrateTables: migrate table columns using GORM
func autoMigrateTables() {
	databases.App.DBConfig.AutoMigrate(&models.Article{})
}

package main

import (
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/controllers"
	"os"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	app := controllers.App{}

	//HOST := os.Getenv("DB_HOST")
	//PORT := os.Getenv("DB_PORT")
	//USER := os.Getenv("DB_USER")
	//NAME := os.Getenv("DB_USER")
	//PASS := os.Getenv("DB_PASS")
	//app.Initialize(HOST, PORT, USER, NAME, PASS)

	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))

	app.RunServer()
}

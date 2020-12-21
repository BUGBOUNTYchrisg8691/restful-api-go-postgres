package controllers

import (
	"fmt"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/middlewares"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/models"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/responses"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
	"log"
	"net/http"
)

type App struct {
	Router 	*mux.Router
	DB		*gorm.DB
}

// Initialize connect to the database and wire up routes
func (app *App) Initialize(DBHost, DBPort, DBUser, DBName, DBPassword string) {
	var err error
	//fmt.Printf("host=%s post=%s user=%s dbname=%s sslmode=disable password=%s",
	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_NAME"),
	//	os.Getenv("DB_PASSWORD"))
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)

	app.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Sprintf("\nCannot connect to database %s", DBName)
		log.Fatal("This is the error: ", err)
	} else {
		fmt.Printf("We are connected to the database %s", DBName)
	}

	//if err != nil {
	//	fmt.Printf("\n Cannot connect to database %s", DBName)
	//	log.Fatal("This is the error: ", err)
	//} else {
	//	fmt.Printf("We are connected to the database %s", DBName)
	//}

	// database migration
	app.DB.Debug().AutoMigrate(&models.User{}, &models.Venue{})

	app.Router = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.Use(middlewares.SetContentTypeMiddleware)	// setting content-type to json

	app.Router.HandleFunc("/", home).Methods("GET")
	app.Router.HandleFunc("/register", app.UserSignUp).Methods("POST")
	app.Router.HandleFunc("/login", app.Login).Methods("POST")

	subrouter := app.Router.PathPrefix("/api").Subrouter()  // routes that require auth
	subrouter.Use(middlewares.AuthJwtVerify)

	subrouter.HandleFunc("/users", app.GetUsers).Methods("GET")
	subrouter.HandleFunc("/venues", app.CreateVenue).Methods("POST")
	subrouter.HandleFunc("/venues", app.GetVenues).Methods("GET")
	subrouter.HandleFunc("/venues/{id:[0-9]+}", app.UpdateVenue).Methods("PUT")
	subrouter.HandleFunc("/venues/{id:[0-9]+}", app.DeleteVenue).Methods("DELETE")
}

func (app *App) RunServer() {
	log.Printf("\nServer starting on port 5000")
	log.Fatal(http.ListenAndServe(":5000", app.Router))
}

func home(writer http.ResponseWriter, request *http.Request) {	// this is the home route
	responses.JSON(writer, http.StatusOK, "Welcome To Ivents")
}
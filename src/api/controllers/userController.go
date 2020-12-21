package controllers

import (
	"encoding/json"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/models"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/responses"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/utils"
	"io/ioutil"
	"net/http"
)

// UserSignUp controller for creating new users
func (app *App) UserSignUp(writer http.ResponseWriter, req *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Registered Successfully"}

	user := &models.User{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(app.DB)
	if usr != nil {
		resp["status"] = "failed"
		resp["message"] = "User already registered, please login"
		responses.JSON(writer, http.StatusBadRequest, resp)
		return
	}

	user.Prepare()	// here strip the text of whitespaces

	err = user.Validate("")	// default were all fields(email, firstname, lastname, password, profileImage) are validated
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
	}

	userCreated, err := user.SaveUser(app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	resp["user"] = userCreated
	responses.JSON(writer, http.StatusCreated, resp)
	return
}

// Login signs in users
func (app *App) Login(writer http.ResponseWriter, req *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Logged In"}

	user := &models.User{}
	body, err := ioutil.ReadAll(req.Body)	// Read user input from request
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	user.Prepare()	// here were strip the text of all whitespace

	err = user.Validate("login")	// fields(email, password) are validated
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUser(app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	if usr == nil {	//user is not registered
		resp["status"] = "failed"
		resp["message"] = "Login failed, please signup"
		responses.JSON(writer, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed, please try again"
		responses.JSON(writer, http.StatusBadRequest, resp)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}
	
	resp["token"] = token
	responses.JSON(writer, http.StatusOK, resp)
	return
}

func (app *App) GetUsers(writer http.ResponseWriter, req *http.Request) {
	users, err := models.GetUsers(app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}
	responses.JSON(writer, http.StatusOK, users)
	return
}
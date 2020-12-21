package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/models"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/responses"
	"github.com/gorilla/mux"
)

// CreateVenue parses request, validates data and saves the new venue
func (app *App) CreateVenue(writer http.ResponseWriter, req *http.Request) {
	var resp = map[string]interface{}{"status": "success",
		"message": "Venue successfully created"}

	user := req.Context().Value("userID").(float64)
	venue := &models.Venue{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &venue)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	venue.Prepare()   // strip whitespace

	if err = venue.Validate(); err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	if vne, _ := venue.GetVenue(app.DB); vne != nil {
		resp["status"] = "failed"
		resp["message"] = "Venue already registered, please choose another name"
		responses.JSON(writer, http.StatusBadRequest, resp)
		return
	}

	venue.UserID = uint(user)

	venueCreated, err := venue.Save(app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	resp["venue"] = venueCreated
	responses.JSON(writer, http.StatusCreated, resp)
	return
}

func (app *App) GetVenues(writer http.ResponseWriter, req *http.Request) {
	venues, err := models.GetVenues(app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}
	responses.JSON(writer, http.StatusOK, venues)
	return
}

func (app *App) UpdateVenue(writer http.ResponseWriter, req *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Venue updated successfully"}

	vars := mux.Vars(req)

	user := req.Context().Value("userID").(float64)
	userID := uint(user)

	id, _ := strconv.Atoi(vars["id"])

	venue, err := models.GetVenueById(id, app.DB)

	if venue.UserID != userID {
		resp["status"] = "failed"
		resp["message"] = "Unauthorized venue update"
		responses.JSON(writer, http.StatusUnauthorized, resp)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	venueUpdate := models.Venue{}
	if err = json.Unmarshal(body, &venueUpdate); err != nil {
		responses.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	venueUpdate.Prepare()

	_, err = venueUpdate.UpdateVenue(id, app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, resp)
	return
}

func (app *App) DeleteVenue(writer http.ResponseWriter, req *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Venue deleted successfully"}

	vars := mux.Vars(req)

	user := req.Context().Value("userID").(float64)
	userID := uint(user)

	id, _ := strconv.Atoi(vars["id"])

	venue, err := models.GetVenueById(id, app.DB)
	if venue.UserID != userID {
		resp["status"] = "failed"
		resp["message"] = "Unauthorized venue delete"
		responses.JSON(writer, http.StatusUnauthorized, resp)
		return
	}

	err = models.DeleteVenue(id, app.DB)
	if err != nil {
		responses.ERROR(writer, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(writer, http.StatusOK, resp)
	return
}
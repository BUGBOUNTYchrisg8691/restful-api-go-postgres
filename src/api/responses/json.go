package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON returns a well formatted response with status code
func JSON(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(data)

	if err != nil {
		fmt.Fprintf(writer, "%s", err.Error())
	}
}

// ERROR return a jsonified error response with status code
func ERROR(writer http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(writer, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(writer, statusCode, nil)
}
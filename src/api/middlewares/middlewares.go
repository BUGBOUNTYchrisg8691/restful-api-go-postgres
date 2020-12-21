package middlewares

import (
	"context"
	"github.com/bugbountychris8691/restful-api-go-postgres/src/api/responses"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

// SetContentTypeMiddleware sets content-type to json
func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, req)
	})
}

// AuthJwtVerify verifies token and add userID to the request context
func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "message":"Missing authorization token"}

		var header = req.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			responses.JSON(writer, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return 	[]byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			resp["status"] = "failed"
			resp["message"] = "Invalid token, please login"
			responses.JSON(writer, http.StatusForbidden, resp)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(req.Context(), "userID", claims["userID"])	// adding the userID to the context
		next.ServeHTTP(writer, req.WithContext(ctx))
	})
}
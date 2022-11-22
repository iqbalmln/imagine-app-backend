package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"gitlab.privy.id/go_graphql/internal/entity"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("fgusbhinjklergoiernglkengkjerbngkerugerb8367yt8734597012840uohrgkdngkerbng")

func JwtAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		claims := &entity.MyClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return JWT_SIGNATURE_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				response := map[string]string{"message": "Unauthorized"}
				ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case jwt.ValidationErrorExpired:
				response := map[string]string{"message": "Unauthorized, Token Expired"}
				ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				response := map[string]string{"message": "Unauthorized sdsdsd"}
				fmt.Println(err)
				ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]string{"message": "Unauthorized"}
			ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		if claims.Role != "Staff IT" {
			response := map[string]string{"message": "Unauthorized as Staff IT"}
			ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		entity.PIC = claims.ID
		next.ServeHTTP(w, r)

	})
}

func ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

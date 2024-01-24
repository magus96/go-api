package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type customClaims struct {
	Number int64 `json:"number"`
	jwt.RegisteredClaims
}

var hmacSampleSecret []byte

func withJWTOut(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Out middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			writeJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return
		}

		if !token.Valid {
			writeJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
		}

		uid, _ := strconv.Atoi(mux.Vars(r)["id"])
		claims := token.Claims.(jwt.MapClaims)
		account, err := s.GetAccountbyID(uid)

		if err != nil {
			writeJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return
		}
		if account.Number != int64(claims["number"].(float64)) {
			writeJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := customClaims{
		account.Number,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "Admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	return ss, err
}

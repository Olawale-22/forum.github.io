package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func CheckUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user's cookie
		cookie, err := r.Cookie("userID")
		if err != nil {
			http.Error(w, "Not logged in", http.StatusUnauthorized)
			return
		}

		// Print the user's activity
		fmt.Printf("User %s made a request to %s\n", cookie.Value, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

/*func ValidCookie(w http.ResponseWriter, r *http.Request) error {
	//var flag bool
	_, err := r.Cookie("userID")
	if err != nil {
		return err
	}
	//flag = true
	return nil
}*/

func ValidCookie(w http.ResponseWriter, r *http.Request) http.Cookie {
	//var flag bool
	cookie, err := r.Cookie("userID")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
	}

	return *cookie
}

// function "GenerateCookie" generates cookie and save it to user browser on login...

func GenerateCookie(w http.ResponseWriter, r *http.Request) http.Cookie {
	// Generate a unique cookie for the user
	expiration := time.Now().Add(365 * 24 * time.Hour)
	//expiration := time.Now().Add(time.Hour)
	id, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cookie := http.Cookie{Name: "userID", Value: id.String(), Expires: expiration, Path: "/"}
	http.SetCookie(w, &cookie)

	return cookie

	// fmt.Fprintln(w, "Cookie generated and saved to browser!")
}

func Itoa(n int) string {
	return strconv.Itoa(n)
}

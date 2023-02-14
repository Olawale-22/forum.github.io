/*

project: forum
SQLite Go-documentation and package used: https://pkg.go.dev/github.com/mattn/go-sqlite3
BUILD- DOCKER IMAGE: "docker build -t myapp ."
RUN - DOCKER CONTAINER: "docker run -p 8282:8282 myapp"

*/

package main

import (
	gt "clouds/getreg"
	hub "clouds/github"
	"log"
	"net/http"
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	http.HandleFunc("/", gt.IndexHandler)
	http.HandleFunc("/indexlog", gt.Homehandler)
	http.HandleFunc("/sign-out", gt.SignOut)
	http.HandleFunc("/forms-profile", gt.ProfileForms)
	http.HandleFunc("/page-register", gt.FillForm)
	http.HandleFunc("/post", gt.PostHandler)
	http.HandleFunc("/page-login", gt.Login)

	// google Auth Handlers
	http.HandleFunc("/gugu", gt.HandleLogin)
	http.HandleFunc("/callback", gt.HandleCallback)

	// github Auth Handlers
	http.HandleFunc("/githb", hub.GithubLoginHandler)
	http.HandleFunc("/githb/callback", hub.GithubCallbackHandler)

	//fmt.Printf("Starting server for testing HTTP POST on https://localhost:8080 ...\n")
	if err := http.ListenAndServe(":8282", nil); err != nil {
		log.Fatal(err)
	}

	/*if err := http.ListenAndServeTLS("127.0.0.1:8080", "/Users/cloud_roi/localhost.crt", "/Users/cloud_roi/localhost.key", nil); err != nil {
		log.Fatal(err)
	}*/
}

package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type MyJsonName struct {
	Webb Web `json:"web"`
}

type Web struct {
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	AuthURI                 string   `json:"auth_uri"`
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret"`
	JavascriptOrigins       []string `json:"javascript_origins"`
	ProjectID               string   `json:"project_id"`
	RedirectUris            []string `json:"redirect_uris"`
	TokenURI                string   `json:"token_uri"`
}

func MashallXIV() Web {
	file, err := os.Open("./middleware/google.json")

	if err != nil {
		fmt.Println("Error opening JSON file:", err)
	}
	defer file.Close()

	data := MyJsonName{}

	web := Web{}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
	}

	web.AuthProviderX509CertURL = data.Webb.AuthProviderX509CertURL
	web.AuthURI = data.Webb.AuthURI
	web.ClientID = data.Webb.ClientID
	web.ClientSecret = data.Webb.ClientSecret
	web.RedirectUris = data.Webb.RedirectUris

	return web
}

var deebee = "./datab.db"

func Init() {
	conn, err := sql.Open("sqlite3", deebee)
	if err != nil {
		fmt.Println("unable to open database home handler")
		//log.Fatal(err.Error())
	}
	conn.Exec("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, full_name TEXT, education TEXT, home_address TEXT, city TEXT, postal_code TEXT, flat_number TEXT, phone TEXT, email TEXT, password TEXT)")

	defer conn.Close()

	_, errK := conn.Exec("CREATE TABLE IF NOT EXISTS activities (user_id INTEGER, post_id INTEGER PRIMARY KEY AUTOINCREMENT, post TEXT, username TEXT, category TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
	if errK != nil {
		fmt.Println(errK)
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS likes_ (user_id INTEGER, username TEXT, post_id INTEGER,  like_id INTEGER PRIMARY KEY AUTOINCREMENT, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(post_id) REFERENCES activities(post_id))")

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS clikes_ (user_id INTEGER, username TEXT, comment_id INTEGER,  like_id INTEGER PRIMARY KEY AUTOINCREMENT, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(comment_id) REFERENCES comments(comment_id))")

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS dislikes_ (user_id INTEGER, username TEXT, post_id INTEGER,  like_id INTEGER PRIMARY KEY AUTOINCREMENT, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(post_id) REFERENCES activities(post_id))")

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS cdislikes_ (user_id INTEGER, username TEXT, comment_id INTEGER,  like_id INTEGER PRIMARY KEY AUTOINCREMENT, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(comment_id) REFERENCES comments(comment_id))")

	if err != nil {
		log.Fatal(err)
	}

	_, errC := conn.Exec("CREATE TABLE IF NOT EXISTS comments (user_id INTEGER, username TEXT, comment_id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, comment TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(post_id) REFERENCES activities(post_id))")

	if errC != nil {
		fmt.Println(errC)
	}

}

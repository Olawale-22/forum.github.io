package github

import (
	"context"
	"fmt"
	"net/http"

	golog "clouds/getreg"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// func init() {
// 	// loads values from .env into the system
// 	if err := godotenv.Load(".env"); err != nil {
// 		log.Fatal("No .env file found")
// 	}
// }

// /*const (
// 	clientID     = "836b569d-c6b4-4b3d-a64d-5de327e3c378"
// 	clientSecret = "YVpZvGNrn3hQrVzZxtKYd1RSDWYgfk5lYxia0ssUngm0"
// 	redirectURI  = "http://localhost:8282/github/callback"
// )*/

// // func init() {
// // 	// loads values from .env in the same directory as the folder containing the Go file into the system
// // 	root, err := os.Getwd()
// // 	if err != nil {
// // 		log.Fatal("Error getting current working directory")
// // 	}
// // 	parentDir := filepath.Dir(root)
// // 	if err := godotenv.Load(filepath.Join(parentDir, ".env")); err != nil {
// // 		log.Fatal("No .env file found in the same directory as the folder containing the Go file")
// // 	}
// // }

// func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get the environment variable
// 	githubClientID := getGithubClientID()

// 	// Create the dynamic redirect URL for login
// 	redirectURL := fmt.Sprintf(
// 		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
// 		githubClientID,
// 		"http://localhost:8282/github/callback",
// 	)

// 	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
// }

// func getGithubClientID() string {

// 	githubClientID, exists := os.LookupEnv("CLIENT_ID")
// 	if !exists {
// 		log.Fatal("Github Client ID not defined in .env file")
// 	}

// 	return githubClientID
// }

// func getGithubClientSecret() string {

// 	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
// 	if !exists {
// 		log.Fatal("Github Client ID not defined in .env file")
// 	}

// 	return githubClientSecret
// }

// func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
// 	if githubData == "" {
// 		// Unauthorized users get an unauthorized message
// 		fmt.Fprintf(w, "UNAUTHORIZED!")
// 		return
// 	}

// 	// Set return type JSON
// 	w.Header().Set("Content-type", "application/json")

// 	// Prettifying the json
// 	var prettyJSON bytes.Buffer
// 	// json.indent is a library utility function to prettify JSON indentation
// 	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
// 	if parserr != nil {
// 		log.Panic("JSON parse error")
// 	}

// 	// Return the prettified JSON as a string
// 	fmt.Println("PRETTIFIED GITHUB JSON", prettyJSON.String())
// }

// func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")

// 	githubAccessToken := getGithubAccessToken(code)

// 	githubData := getGithubData(githubAccessToken)

// 	loggedinHandler(w, r, githubData)
// }

// //****************

// func getGithubData(accessToken string) string {
// 	// Get request to a set URL
// 	req, reqerr := http.NewRequest(
// 		"GET",
// 		"https://api.github.com/user",
// 		nil,
// 	)
// 	if reqerr != nil {
// 		log.Panic("API Request creation failed")
// 	}

// 	// Set the Authorization header before sending the request
// 	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
// 	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
// 	req.Header.Set("Authorization", authorizationHeaderValue)

// 	// Make the request
// 	resp, resperr := http.DefaultClient.Do(req)
// 	if resperr != nil {
// 		log.Panic("Request failed")
// 	}

// 	// Read the response as a byte slice
// 	respbody, _ := ioutil.ReadAll(resp.Body)

// 	// Convert byte slice to string and return
// 	return string(respbody)
// }

// type githubAccessTokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// 	Scope       string `json:"scope"`
// }

// func getGithubAccessToken(code string) string {

// 	clientID := getGithubClientID()
// 	clientSecret := getGithubClientSecret()

// 	// Set us the request body as JSON
// 	requestBodyMap := map[string]string{
// 		"client_id":     clientID,
// 		"client_secret": clientSecret,
// 		"code":          code,
// 	}
// 	requestJSON, _ := json.Marshal(requestBodyMap)

// 	// POST request to set URL
// 	req, reqerr := http.NewRequest(
// 		"POST",
// 		"https://github.com/login/oauth/access_token",
// 		bytes.NewBuffer(requestJSON),
// 	)

// 	if reqerr != nil {
// 		log.Panic("Request creation failed")
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Accept", "application/json")

// 	// Get the response
// 	resp, resperr := http.DefaultClient.Do(req)
// 	if resperr != nil {
// 		log.Panic("Request failed")
// 	}

// 	// Response body converted to stringified JSON
// 	respbody, _ := ioutil.ReadAll(resp.Body)

// 	// Represents the response received from Github

// 	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
// 	var ghresp githubAccessTokenResponse
// 	json.Unmarshal(respbody, &ghresp)

// 	// Return the access token (as the rest of the
// 	// details are relatively unnecessary for us)
// 	return ghresp.AccessToken
// }

const (
	clientID     = "22cc907b9cc7adacb953"
	clientSecret = "1d2a1a105cb9b1888ce5826f7bd5880e26c74864"
	redirectURI  = "http://localhost:8282/githb/callback"
)

var (
	conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		Scopes: []string{"user:email"},
	}
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", clientID, redirectURI)
	http.Redirect(w, r, url, http.StatusFound)
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// Exchange authorization code for access token
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprint(w, "Error exchanging code for token")
		return
	}

	// Create GitHub client using access token
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	client := github.NewClient(tc)

	// Get authenticated user's email
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Fprint(w, "Error getting user info")
		return
	}

	email, _, err := client.Users.ListEmails(context.Background(), nil)
	if err != nil {
		fmt.Fprint(w, "Error getting user email")
		return
	}

	username := fmt.Sprintln(*user.Login)
	mail := fmt.Sprintln(email[0].GetEmail())
	password := fmt.Sprintln(*user.ID)

	golog.GoLogGoogle(w, mail, username, password, r)

	//fmt.Printf("----<<<<>>>>>>>>>>>Welcome, %s</h1><p>Your email is %s and ID is %v <>----<<<<<<<<<<>>>>>>>", *user.Login, email[0].GetEmail(), *user.ID)
	//GrantAccess(w, "/indexlog", r)

}

/*func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}*/

/*
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", clientID, redirectURI)
	http.Redirect(w, r, url, http.StatusFound)
}
*/

func GrantAccess(w http.ResponseWriter, s string, r *http.Request) {
	http.Redirect(w, r, s, http.StatusFound)
	//http.Redirect(w, r, s, http.StatusSeeOther)
}

package getreg

//****** PROJECT NAME: Zone_Forum ******
//****** DEVELOPER: Olawale22 Sulaiaman ******

import (
	"bytes"
	autth "clouds/autth"
	cookii "clouds/middleware"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Registration struct {
	User     string `json:"username"`
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type RegistrationDB struct {
	New []*Registration
}

var DB = &RegistrationDB{}

func (u *RegistrationDB) addUser(reg Registration) {
	u.New = append(u.New, &reg)
}

func notExist(u Registration) bool {
	//ok := true
	for _, v := range DB.New {
		if (u.User == v.User || u.Mail == v.Mail) && (u.Password == v.Password) {
			return false
		}
	}
	return true
}

var Deebee = "./datab.db"

type Posts struct {
	P_Uid          int
	Pid            int
	Name           string
	Category       string
	Posted         string
	Likers         []string
	Likes          int
	Dislikes       int
	Dislikers      []string
	ScrollPosition string
	Comment        []Comments
	Img            string
	Time           string
	Commenters     []Comments
}

type Comments struct {
	Id             int
	Pid            int
	Cid            int
	Name           string
	Comment        string
	Img            string
	Time           string
	ScrollPosition int
	Likes          int
	Dislikes       int
	//Time    time.Time

}

type User struct {
	ID          int
	Name        string
	Img         string
	FullName    string
	Education   string
	HomeAddress string
	Phone       string
}

type SDK struct {
	Us  User
	Pst []Posts
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return
		}

		conn, err := sql.Open("sqlite3", Deebee)
		if err != nil {
			fmt.Println("unable to open database home handler")
			//log.Fatal(err.Error())
		}

		defer conn.Close()

		// _, errK := conn.Exec("CREATE TABLE IF NOT EXISTS activities (user_id INTEGER, post_id INTEGER PRIMARY KEY AUTOINCREMENT, post TEXT, username TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
		// if errK != nil {
		// 	fmt.Println(errK)
		// }

		var sdk []Posts
		sth, errNext := conn.Query("SELECT user_id, post_id, username, post, STRFTIME('%H:%M', date_created) FROM activities")
		if errNext != nil {
			fmt.Println(errNext)
			return
		}

		defer sth.Close()

		for sth.Next() {
			var str Posts
			//str.Time = time.Now()
			err = sth.Scan(&str.P_Uid, &str.Pid, &str.Name, &str.Posted, &str.Time)
			if err != nil {
				fmt.Println(err)
				return
			}
			sdk = append(sdk, str)
		}

		var newSdk []Posts
		for i := len(sdk) - 1; i > 0; i-- {
			newSdk = append(newSdk, sdk[i])
		}

		fmt.Println("ROWS HERE ->", sdk)
		fmt.Println(w, "Form data saved successfully!")
		//Maketmpl(w, "index", nil)
		Maketmpl(w, "index", newSdk)
		return
	} else if r.Method == http.MethodPost {
		//fmt.Fprintln(w, "Sign In")
		GrantAccess(w, "/page-login", r)
		return
	}
}

func ProfileForms(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		Maketmpl(w, "forms-profile", nil)

	} else if r.Method == http.MethodPost {

		fullName := r.FormValue("fullName")
		education := r.FormValue("education")
		address := r.FormValue("streetAddress")
		city := r.FormValue("city")
		postalCode := r.FormValue("postalCode")
		Appartment := r.FormValue("apt")
		phone := r.FormValue("phone")

		conn, err := sql.Open("sqlite3", Deebee)
		if err != nil {
			fmt.Println(err)
		}

		cookies, err := r.Cookie("userID")

		if err != nil {
			fmt.Println("Homehandler cookie not valid")
			fmt.Println(err)
			return
		}

		//var usingAname string
		var hydee int
		errNew := conn.QueryRow("SELECT user_id FROM logs WHERE coookies = ?", cookies.Value).Scan(&hydee)
		if errNew != nil {
			fmt.Println(errNew)
			return
		}

		home := address + ", flat " + Appartment
		_, err = conn.Exec("UPDATE users SET full_name=?, education=?, home_address=?, city=?, postal_code=?, flat_number=?, phone=? WHERE user_id=?", fullName, education, home, city, postalCode, Appartment, phone, hydee)
		if err != nil {
			log.Fatal(err)
		}

		GrantAccess(w, "/indexlog", r)
		//http.Redirect(w, r, "/indexlog", http.StatusFound)
		return
		//conn.Exec("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, full_name TEXT, education TEXT, home_address TEXT, city TEXT, postal_code TEXT, flat_number TEXT, phone TEXT, email TEXT, password TEXT)")

	}
}

//********

func SignOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "userID",
		Value:  "",
		Path:   "/indexlog",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	// Redirect to the login page
	http.Redirect(w, r, "/page-login", http.StatusFound)
	//GrantAccess(w, "/page-login", r)
}

type Liks struct {
	LikeName string
	Postid   int
	UserId   int
}

//***************

func redeemPassword(password string) bool {
	pattern := `^([A-Z]*?[0-9]*?[#\$-€]*?).{8,}$`
	match, _ := regexp.MatchString(pattern, password)
	return match
}

func DBver(conn *sql.DB, a, b string) bool {
	var exists bool
	errNew := conn.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND username = ?", a, b).Scan(&exists)
	if errNew != nil {
		fmt.Println(errNew)
	}
	return exists
}

func Itoa(n int) string {
	return strconv.Itoa(n)
}

func FillForm(w http.ResponseWriter, r *http.Request) {
	//res := mux.NewRouter()
	if r.Method == http.MethodGet {
		Maketmpl(w, "page-register", nil)

	} else if r.Method == http.MethodPost {
		var count int
		conn, err := sql.Open("sqlite3", Deebee)
		if err != nil {
			//fmt.Println("unable to open database")
			log.Fatal(err.Error())
		}

		defer conn.Close()

		// conn.Exec("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, full_name TEXT, education TEXT, home_address TEXT, city TEXT, postal_code TEXT, flat_number TEXT, phone TEXT, email TEXT, password TEXT)")

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uName := r.FormValue("uname")
		eMail := r.FormValue("eemail")
		pWord := r.FormValue("pwd")

		errNew := conn.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ?", eMail, uName).Scan(&count)
		if errNew != nil {
			fmt.Println(errNew)
			return
		}

		if count > 0 {
			fmt.Fprintln(w, "Email or username already in use.")
			return

		} else if !redeemPassword(pWord) {
			fmt.Fprintln(w, "password should have at least one capital letter, one number and one or more #$-€")
			return
		} else {

			input, _ := conn.Prepare("INSERT INTO users (username, email, password) VALUES(?, ?, ?)")
			input.Exec(uName, eMail, EncryptPassword(pWord))
			fmt.Println("Form data saved successfully!")

			//check input in database
			rows, _ := conn.Query("SELECT user_id, username, email, password FROM users")

			//var id int
			var username string
			var email string
			var password string
			var id int

			// output database on the terminal "rows"
			for rows.Next() {
				rows.Scan(&id, &username, &email, &password)
				fmt.Println("user_id:", id, ", name:, ", username, "email:", email, ", password:", password)
				fmt.Println("pwd before encryption:", pWord)
			}

			inpt := Registration{
				User:     uName,
				Mail:     eMail,
				Password: pWord,
			}

			//DB := &RegistrationDB{}
			ok := notExist(inpt)
			if !ok {
				fmt.Fprintf(w, "Invalid login")
			} else {
				DB.addUser(inpt)
				// Convert the User object to JSON
				/*userJSON, err := json.Marshal(DB)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}*/
				userJSON := json.NewEncoder(os.Stdout)
				userJSON.SetIndent("", "  ")
				if err := userJSON.Encode(&DB); err != nil {
					panic(err)
				}
			}
			GrantAccess(w, "/page-login", r)
		}
	}
}

func EncryptPassword(s string) string {
	res := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x\n", res)
}

type LogInfo struct {
	Name  string
	Email string
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Maketmpl(w, "page-login", nil)
	} else if r.Method == http.MethodPost {
		conn, err := sql.Open("sqlite3", Deebee)
		if err != nil {
			//fmt.Println("unable to open database")
			log.Fatal(err.Error())
		}

		defer conn.Close()
		//conn.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, coookie BLOB, FOREIGN KEY(id) REFERENCES users(user_id))")

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userName := r.FormValue("logmail")
		//eMail := r.FormValue("eemail")
		passWord := r.FormValue("logpwd")

		//display errData info when wrong username or password is provided by user
		errData := struct {
			ErrUsername string
			ErrPassword string
		}{
			ErrUsername: "Yo wtf i dont know you What's your Username",
			ErrPassword: "Username or Password unknown",
		}

		var dbUsername, dbPassword string
		var dbFullName, dbEdu, dbHomeAddress, dbPhone sql.NullString
		var out_id int

		//var counter int

		// conn.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", userName).Scan(&counter)
		// // if errJ != nil {

		// // 	if err == sql.ErrNoRows {
		// // 		errData.ErrPassword = ""
		// // 		Maketmpl(w, "page-login", errData)
		// // 	}
		// // }

		// if counter <= 0 {
		// 	errData.ErrPassword = ""
		// 	Maketmpl(w, "page-login", errData)
		// } else {

		errNN := conn.QueryRow("SELECT user_id, username, password, full_name, education, home_address, phone FROM users WHERE username = ?", userName).Scan(&out_id, &dbUsername, &dbPassword, &dbFullName, &dbEdu, &dbHomeAddress, &dbPhone)

		if errNN != nil {

			if err == sql.ErrNoRows {
				errData.ErrPassword = ""
				Maketmpl(w, "page-login", errData)
				return
			}
			//log.Fatal(errNN)
		}

		// check if the provided "password" matches the one in the "users" table
		if EncryptPassword(passWord) != dbPassword {
			errData.ErrUsername = ""
			Maketmpl(w, "page-login", errData)
			return
		} else {

			// _, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
			// if err != nil {
			// 	log.Fatal(err)
			// }

			/*if !redeemPassword(passWord) {
				fmt.Fprintln(w, "password should have at least one capital letter, one number and one or more #$-€")
				return

			}*/

			ck := cookii.GenerateCookie(w, r)
			_, err = conn.Exec("INSERT INTO logs (user_id, username, coookies, timestamp) VALUES (?,?,?, datetime('now'))", out_id, dbUsername, ck.Value)
			if err != nil {
				log.Fatal(err)
			}

			//***********************

			fmt.Println("cookie.Value: ", ck)
			//input, _ := conn.Prepare("INSERT INTO logs (user_id, coookie) VALUES(?, ?)")
			//input.Exec(out_id, ck)
			fmt.Printf("New sign in: %s %s", userName, passWord)
			//GrantAccess(w, "/indexlog", r)
			if !dbFullName.Valid && !dbEdu.Valid && !dbPhone.Valid {
				GrantAccess(w, "/forms-profile", r)
				return

			} else {
				GrantAccess(w, "/indexlog", r)
				return
			}
		}
	}
}

func GrantAccess(w http.ResponseWriter, s string, r *http.Request) {
	http.Redirect(w, r, s, http.StatusFound)
	//http.Redirect(w, r, s, http.StatusSeeOther)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Maketmpl(w http.ResponseWriter, tmplName string, data interface{}) {

	templateCache, err := createTemplateCache()

	if err != nil {
		checkErr(err)
	}

	tpl, err2d2 := templateCache[tmplName+".html"]
	if !err2d2 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buff := new(bytes.Buffer)
	tpl.Execute(buff, data)
	buff.WriteTo(w)
}

func createTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./static/*.html")
	if err != nil {
		return cache, nil
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))
		if err != nil {
			return cache, nil
		}
		cache[name] = tmpl
	}
	return cache, nil
}

func Atoi(s string) int {
	ok, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	return ok
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
}

var (
	ask               = autth.MashallXIV()
	googleOauthConfig = &oauth2.Config{

		RedirectURL:  ask.RedirectUris[0],
		ClientID:     ask.ClientID,
		ClientSecret: ask.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "zone_forum"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

// WORKS PERFECT
func HandleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != oauthStateString {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	ctx := context.Background()
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get user information from the ID token in the token response

	// For example, you can use the Google Userinfo API to get user information:

	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "failed to get user info: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read user info: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("******************", string(contents))
	var userInfo UserInfo
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		fmt.Println("failed to unmarshal user info:", err)
		return
	}

	email := userInfo.Email
	pass := userInfo.Sub
	userN := userInfo.GivenName

	GoLogGoogle(w, email, userN, pass, r)

}

func GoLogGoogle(w http.ResponseWriter, email string, username string, password string, r *http.Request) {

	conn, err := sql.Open("sqlite3", Deebee)
	if err != nil {
		fmt.Println("unable to open database Golog")
		log.Fatal(err.Error())
	}

	defer conn.Close()

	//conn.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, coookie BLOB, FOREIGN KEY(id) REFERENCES users(user_id))")

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUsername, dbPassword string
	var dbFullName, dbEdu, dbHomeAddress, dbPhone sql.NullString
	var out_id int

	counter := 0

	errNN := conn.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&counter)
	if errNN != nil {

		fmt.Println("unable to query users Golog==", errNN)
		return
	}

	if counter > 0 {
		errNN := conn.QueryRow("SELECT user_id, username, password, full_name, education, home_address, phone FROM users WHERE username = ?", username).Scan(&out_id, &dbUsername, &dbPassword, &dbFullName, &dbEdu, &dbHomeAddress, &dbPhone)

		if errNN != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			log.Fatal(errNN)
		}

		// check if the provided "password" matches the one in the "users" table
		if EncryptPassword(password) != dbPassword {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		} else {

			ck := cookii.GenerateCookie(w, r)
			_, err = conn.Exec("INSERT INTO logs (user_id, username, coookies, timestamp) VALUES (?,?,?, datetime('now'))", out_id, dbUsername, ck.Value)
			if err != nil {
				fmt.Println("Unable to insert into logs Golog")
				log.Fatal(err)
			}

			fmt.Println("cookie.Value: ", ck)

			//input, _ := conn.Prepare("INSERT INTO logs (user_id, coookie) VALUES(?, ?)")
			//input.Exec(out_id, ck)

			fmt.Printf("New sign in: %s %s", username, password)
			//GrantAccess(w, "/indexlog", r)
			if !dbFullName.Valid && !dbEdu.Valid && !dbPhone.Valid {
				GrantAccess(w, "/forms-profile", r)
				return
			} else {

				GrantAccess(w, "/indexlog", r)
				return
			}
		}

	} else {

		input, _ := conn.Prepare("INSERT INTO users (username, email, password) VALUES(?, ?, ?)")
		input.Exec(username, email, EncryptPassword(password))
		fmt.Println("Form data saved successfully!")
		errNN := conn.QueryRow("SELECT user_id, username, password, full_name, education, home_address, phone FROM users WHERE username = ?", username).Scan(&out_id, &dbUsername, &dbPassword, &dbFullName, &dbEdu, &dbHomeAddress, &dbPhone)

		if errNN != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			fmt.Println("error querying Rows users Gologs")
			log.Fatal(errNN)
		}

		// check if the provided "password" matches the one in the "users" table

		if EncryptPassword(password) != dbPassword {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		} else {

			// _, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
			// if err != nil {
			// 	fmt.Println("Unable to insert into logs Golog")
			// 	log.Fatal(err)
			// }

			ck := cookii.GenerateCookie(w, r)
			_, err = conn.Exec("INSERT INTO logs (user_id, username, coookies, timestamp) VALUES (?,?,?, datetime('now'))", out_id, dbUsername, ck.Value)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("cookie.Value: ", ck)

			fmt.Printf("New sign in: %s %s", username, password)

			//GrantAccess(w, "/indexlog", r)
			if !dbFullName.Valid && !dbEdu.Valid && !dbPhone.Valid {
				GrantAccess(w, "/forms-profile", r)
				return

			} else {
				GrantAccess(w, "/indexlog", r)
				return
			}
		}
	}
}

const (
	clientID     = "22cc907b9cc7adacb953"
	clientSecret = "1d2a1a105cb9b1888ce5826f7bd5880e26c74864"
	redirectURI  = "http://localhost:8282/githb/callback"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", clientID, redirectURI)
	http.Redirect(w, r, url, http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	fmt.Fprintf(w, "Welcome back, your code is: %s", code)
}

// func LikeHandler(w http.ResponseWriter, r *http.Request) {

// 	conn, err := sql.Open("sqlite3", Deebee)
// 	if err != nil {
// 		fmt.Println("unable to open database Golog")
// 		log.Fatal(err.Error())
// 	}

// 	defer conn.Close()

// 	cookies, err := r.Cookie("userID")

// 	if err != nil {
// 		fmt.Println("Homehandler cookie not valid")
// 		fmt.Println(err)
// 		return
// 	}

// 	var usingname string
// 	var hydee int
// 	errNew := conn.QueryRow("SELECT user_id FROM logs WHERE coookies = ?", cookies.Value).Scan(&hydee)

// 	if errNew != nil {
// 		fmt.Println(errNew)
// 		return
// 	}

// 	Pd := r.FormValue("lyks")

// 	errNN := conn.QueryRow("SELECT username FROM likes_ WHERE user_id = ? AND post_id = ?", hydee, Pd).Scan(&usingname)

// 	if errNN != nil {
// 		if err == sql.ErrNoRows {
// 			_, err = conn.Exec("INSERT INTO likes_ (user_id, username, post_id) VALUES (?,?,?)", hydee, usingname, Pd)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			GrantAccess(w, "/indexlog", r)
// 			return
// 		}
// 		fmt.Println("error querying Rows users Gologs")
// 		log.Fatal(errNN)
// 	}
// 	GrantAccess(w, "/indexlog", r)
// 	return
// }

// using golang how can I listen to a button o'clock event and send it to my go backend example <button name="likes", type="submit></button> if I get the button using r.FormValue("likes") I want to know when user clicks this button so that the like can be added to my sqlite3 "likes" table in the database

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
	"time"

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

var deebee = "./datab.db"

type Posts struct {
	P_Uid          int
	Pid            int
	Name           string
	Posted         string
	ScrollPosition string
	Comment        []Comments
	Img            string
	Time           string
	Commenters     []Comments
}

type Comments struct {
	Id             int
	Pid            int
	Name           string
	Comment        string
	Img            string
	Time           string
	ScrollPosition int
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
		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			fmt.Println("unable to open database home handler")
			//log.Fatal(err.Error())
		}

		defer conn.Close()

		_, errK := conn.Exec("CREATE TABLE IF NOT EXISTS activities (user_id INTEGER, post_id INTEGER PRIMARY KEY AUTOINCREMENT, post TEXT, username TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
		if errK != nil {
			fmt.Println(errK)
		}

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

//****** THIS IS WHERE WE ARE func ProfileForms() ******
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

		conn, err := sql.Open("sqlite3", deebee)
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

func Homehandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Del("Content-Security-Policy")
		w.Header().Set("Cache-Control", "no-cache")

		if r.URL.Path != "/indexlog" {
			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return
		}

		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			fmt.Println("unable to open database home handler")
			//log.Fatal(err.Error())
		}

		defer conn.Close()

		var sdk []Posts
		sth, errNext := conn.Query("SELECT user_id, post_id, username, post, STRFTIME('%H:%M', date_created) FROM activities")
		if errNext != nil {
			fmt.Println(errNext)
			return
		}

		//********************

		cookies, err := r.Cookie("userID")

		if err != nil {
			fmt.Println("Homehandler cookie not valid")
			fmt.Println(err)
			return
		}

		var usingAname string
		var hydee int
		errNew := conn.QueryRow("SELECT user_id FROM logs WHERE coookies = ?", cookies.Value).Scan(&hydee)
		if errNew != nil {
			fmt.Println(errNew)
			return
		}

		var fname, edu, homeAdd, phoone string

		errNexxt := conn.QueryRow("SELECT username, full_name, education, home_address, phone FROM users WHERE user_id = ?", hydee).Scan(&usingAname, &fname, &edu, &homeAdd, &phoone)
		if errNexxt != nil {
			fmt.Println(errNexxt)
			return
		}

		//********************

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

		stN, eer := conn.Query("SELECT user_id, username, post_id, comment, STRFTIME('%H:%M', date_created) FROM comments")
		if eer != nil {
			log.Fatal(eer)
		}

		var ccc Comments

		defer stN.Close()

		for stN.Next() {

			eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Pid, &ccc.Comment, &ccc.Time)

			if eer != nil {
				fmt.Println(eer)
				return
			}

			for v := 0; v < len(sdk); v++ {

				if sdk[v].Pid == ccc.Pid {

					sdk[v].Commenters = append(sdk[v].Commenters, ccc)
					if len(sdk[v].Commenters) > 8 {
						sdk[v].Commenters = sdk[v].Commenters[0:9]

					}

				}

			}
		}

		var newSdk SDK
		newSdk.Us.ID = hydee
		newSdk.Us.Name = usingAname
		newSdk.Us.Education = edu
		newSdk.Us.FullName = fname
		newSdk.Us.HomeAddress = homeAdd
		newSdk.Us.Phone = phoone

		for i := len(sdk) - 1; i > 0; i-- {
			newSdk.Pst = append(newSdk.Pst, sdk[i])
		}

		fmt.Println("ROWS HERE ->", sdk)
		fmt.Println(w, "Form data saved successfully!")
		Maketmpl(w, "indexlog", newSdk)

		//Maketmpl(w, "indexlog", sdk)
	} else if r.Method == http.MethodPost {

		w.Header().Set("Cache-Control", "no-cache")
		//w.Write([]byte("POST request processed"))
		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			fmt.Println("unable to open database home handler")
			//log.Fatal(err.Error())
		}

		defer conn.Close()

		cookies, err := r.Cookie("userID")

		if err != nil {
			fmt.Println("Homehandler cookie not valid")
			fmt.Println(err)
			return
		}

		var usingAname string
		var hydee int
		errNew := conn.QueryRow("SELECT user_id FROM logs WHERE coookies = ?", cookies.Value).Scan(&hydee)

		if errNew != nil {
			fmt.Println(errNew)
			return
		}
		//var fname, edu, homeAdd, phoone sql.NullString

		var fname, edu, homeAdd, phoone string

		errNexxt := conn.QueryRow("SELECT username, full_name, education, home_address, phone FROM users WHERE user_id = ?", hydee).Scan(&usingAname, &fname, &edu, &homeAdd, &phoone)

		if errNexxt != nil {
			fmt.Println(errNexxt)
			return
		}

		//var str Posts
		//var Allcmt []Comments

		_, errK := conn.Exec("CREATE TABLE IF NOT EXISTS comments (user_id INTEGER, username TEXT, comment_id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, comment TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(post_id) REFERENCES activities(post_id))")

		if errK != nil {
			fmt.Println(errK)
		}

		var sdk []Posts

		posts := r.FormValue("postit")

		Chat := r.FormValue("chat")

		PostId := r.FormValue("post_IId")

		if posts != "" && Chat == "" {

			/****
			_, err = conn.Exec("UPDATE activities SET user_id=?, post=?, username=?, date_created=?", hydee, posts, usingAname, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			****/

			_, err = conn.Exec("INSERT INTO activities (user_id,  post, username, date_created) VALUES(?, ?, ?, datetime('now'))", hydee, posts, usingAname, time.Now().Add(time.Hour))
			if err != nil {
				log.Fatal(err)
			}

			sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities")

			if errNext != nil {
				fmt.Println(errNext)
				return
			}

			defer sth.Close()

			for sth.Next() {

				var str Posts

				//str.Time = time.Now()
				//str.Commenters = Allcmt
				err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
				if err != nil {
					fmt.Println(err)
					return
				}
				sdk = append(sdk, str)
			}

			stN, eer := conn.Query("SELECT user_id, username, post_id, comment, STRFTIME('%H:%M', date_created) FROM comments")
			if eer != nil {
				log.Fatal(eer)
			}

			var ccc Comments

			defer stN.Close()

			for stN.Next() {

				eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Pid, &ccc.Comment, &ccc.Time)

				if eer != nil {
					fmt.Println(eer)
					return
				}

				for v := 0; v < len(sdk); v++ {

					if sdk[v].Pid == ccc.Pid {

						sdk[v].Commenters = append(sdk[v].Commenters, ccc)
						if len(sdk[v].Commenters) > 8 {
							sdk[v].Commenters = sdk[v].Commenters[0:9]

						}

					}

				}
			}

		} else if Chat != "" && posts == "" {
			//scrollPosition := r.FormValue("scrollPosition")

			_, err = conn.Exec("INSERT INTO comments (user_id,  username, post_id, comment, date_created) VALUES(?, ?, ?, ?, datetime('now'))", hydee, usingAname, Atoi(PostId), Chat, time.Now().Add(time.Hour))
			if err != nil {
				log.Fatal(err)
			}

			sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities")

			if errNext != nil {
				fmt.Println(errNext)
				return
			}

			defer sth.Close()

			for sth.Next() {

				var str Posts

				//str.Time = time.Now()
				//str.Commenters = Allcmt

				err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
				if err != nil {
					fmt.Println(err)
					return
				}
				sdk = append(sdk, str)
			}

			/*var cName, cCmt, cTime string
			var cPid int*/

			// add comments to posts
			stN, eer := conn.Query("SELECT user_id, username, post_id, comment, STRFTIME('%H:%M', date_created) FROM comments")
			if eer != nil {
				log.Fatal(eer)
			}

			var ccc Comments

			defer stN.Close()

			for stN.Next() {

				eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Pid, &ccc.Comment, &ccc.Time)
				if eer != nil {
					fmt.Println(eer)
					return
				}

				for v := 0; v < len(sdk); v++ {
					if sdk[v].Pid == ccc.Pid {
						//ccc.ScrollPosition = Atoi(scrollPosition)
						sdk[v].Commenters = append(sdk[v].Commenters, ccc)
						if len(sdk[v].Commenters) > 8 {
							sdk[v].Commenters = sdk[v].Commenters[0:9]
						}
					}

				}
			}

		}

		var newSdk SDK
		newSdk.Us.ID = hydee
		newSdk.Us.Name = usingAname
		newSdk.Us.Education = edu
		newSdk.Us.FullName = fname
		newSdk.Us.HomeAddress = homeAdd
		newSdk.Us.Phone = phoone

		for i := len(sdk) - 1; i > 0; i-- {
			newSdk.Pst = append(newSdk.Pst, sdk[i])
		}

		fmt.Println("ROWS HERE ->", sdk)
		fmt.Println(w, "Form data saved successfully!")
		fmt.Fprintf(w, "<html><head><script>window.scrollTo(0, %s);</script></head></html>", PostId)
		Maketmpl(w, "indexlog", newSdk)
		return
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		ID := r.FormValue("id")

		iD := Atoi(ID)
		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			fmt.Println(err)
			return
		}

		var str Posts
		errNext := conn.QueryRow("SELECT user_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE post_id = ?", iD).Scan(&str.P_Uid, &str.Posted, &str.Name, &str.Time)

		if errNext != nil {
			fmt.Println(errNext)
			return
		}

		stN, eer := conn.Query("SELECT user_id, username, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)
		if eer != nil {
			log.Fatal(eer)
		}

		var ccc Comments
		defer stN.Close()
		for stN.Next() {
			eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Comment, &ccc.Time)
			if eer != nil {

				fmt.Println(eer)
				return

			}

			str.Commenters = append(str.Commenters, ccc)
		}

		fmt.Println("*************** POST PAGE HERE ->", str)
		fmt.Println(w, "Form data saved successfully!")
		Maketmpl(w, "post", str)
		return
	} else if r.Method == http.MethodPost {
		w.Header().Set("Cache-Control", "no-cache")

		ID := r.FormValue("id")

		chats := r.FormValue("chats")

		iD := Atoi(ID)

		if _, err := strconv.Atoi(ID); err != nil {

			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return

		}

		if r.URL.Path != "/post" || iD < 0 {

			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return

		}

		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			fmt.Println(err)
			return
		}

		var str Posts
		errNext := conn.QueryRow("SELECT user_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE post_id = ?", iD).Scan(&str.P_Uid, &str.Posted, &str.Name, &str.Time)

		if errNext != nil {
			fmt.Println(errNext)
			return
		}

		_, err = conn.Exec("INSERT INTO comments (user_id,  username, post_id, comment, date_created) VALUES(?, ?, ?, ?, datetime('now'))", str.P_Uid, &str.Name, iD, chats, time.Now().Add(time.Hour))
		if err != nil {
			log.Fatal(err)
		}

		stN, eer := conn.Query("SELECT user_id, username, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

		if eer != nil {

			log.Fatal(eer)

		}
		var ccc Comments
		defer stN.Close()
		for stN.Next() {

			eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Comment, &ccc.Time)

			if eer != nil {
				fmt.Println(eer)
				return
			}
			str.Commenters = append(str.Commenters, ccc)
		}

		fmt.Println("*************** POST PAGE HERE ->", str)
		fmt.Println(w, "Form data saved successfully!")
		Maketmpl(w, "post", str)
		return
	}
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
		conn, err := sql.Open("sqlite3", deebee)
		if err != nil {
			//fmt.Println("unable to open database")
			log.Fatal(err.Error())
		}

		defer conn.Close()

		conn.Exec("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, full_name TEXT, education TEXT, home_address TEXT, city TEXT, postal_code TEXT, flat_number TEXT, phone TEXT, email TEXT, password TEXT)")

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

		conn, err := sql.Open("sqlite3", deebee)
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

		var dbUsername, dbPassword string
		var dbFullName, dbEdu, dbHomeAddress, dbPhone sql.NullString
		var out_id int

		errNN := conn.QueryRow("SELECT user_id, username, password, full_name, education, home_address, phone FROM users WHERE username = ?", userName).Scan(&out_id, &dbUsername, &dbPassword, &dbFullName, &dbEdu, &dbHomeAddress, &dbPhone)

		if errNN != nil {

			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			log.Fatal(errNN)
		}

		// check if the provided "password" matches the one in the "users" table
		if EncryptPassword(passWord) != dbPassword {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		} else {
			_, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
			if err != nil {
				log.Fatal(err)
			}

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
		panic(err)
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
	//
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

	conn, err := sql.Open("sqlite3", deebee)
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
			_, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
			if err != nil {
				log.Fatal(err)
			}

			/*if !redeemPassword(passWord) {
				fmt.Fprintln(w, "password should have at least one capital letter, one number and one or more #$-€")
				return

			}*/
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
			_, err = conn.Exec("CREATE TABLE IF NOT EXISTS logs (user_id INTEGER, username TEXT, coookies TEXT, timestamp DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id))")
			if err != nil {
				fmt.Println("Unable to insert into logs Golog")
				log.Fatal(err)
			}

			ck := cookii.GenerateCookie(w, r)
			_, err = conn.Exec("INSERT INTO logs (user_id, username, coookies, timestamp) VALUES (?,?,?, datetime('now'))", out_id, dbUsername, ck.Value)
			if err != nil {
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

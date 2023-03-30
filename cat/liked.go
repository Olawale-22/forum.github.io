package cat

import (
	reg "clouds/getreg"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type LIKEE struct {
	Pid int
}

//var DBEE = "./datab.db"

func LikedPosts(w http.ResponseWriter, r *http.Request) {

	//if r.Method == http.MethodGet {
	w.Header().Del("Content-Security-Policy")
	w.Header().Set("Cache-Control", "no-cache")

	if r.URL.Path != "/liked" {
		http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
		return
	}

	conn, err := sql.Open("sqlite3", DBEE)
	if err != nil {
		fmt.Println("unable to open database home handler")
		//log.Fatal(err.Error())
	}

	defer conn.Close()

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
	//var cty int
	stlikes, errNN := conn.Query("SELECT post_id FROM likes_ WHERE user_id = ?", hydee)

	if errNN != nil {
		fmt.Println(errNN)
	}

	defer stlikes.Close()

	var ouu []int
	for stlikes.Next() {
		var stlyk LIKEE
		//str.Time = time.Now()
		err = stlikes.Scan(&stlyk.Pid)
		ouu = append(ouu, stlyk.Pid)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var sdk []reg.Posts
	sth, errNext := conn.Query("SELECT post_id, post, username, category, STRFTIME('%H:%M', date_created) FROM activities")

	if errNext != nil {
		fmt.Println(errNext)
		return
	}

	var fname, edu, homeAdd, phoone string

	errNexxt := conn.QueryRow("SELECT username, full_name, education, home_address, phone FROM users WHERE user_id = ?", hydee).Scan(&usingAname, &fname, &edu, &homeAdd, &phoone)
	if errNexxt != nil {
		fmt.Println(errNexxt)
		return
	}

	//*******************

	defer sth.Close()

	for sth.Next() {

		var str reg.Posts
		//str.Time = time.Now()
		err = sth.Scan(&str.Pid, &str.Posted, &str.Name, &str.Category, &str.Time)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i < len(ouu); i++ {
			if str.Pid == ouu[i] {
				sdk = append(sdk, str)
			}
		}

	}

	stL, _ := conn.Query("SELECT user_id, username, post_id FROM likes_")
	likes := reg.Liks{}

	for stL.Next() {

		stL.Scan(&likes.UserId, &likes.LikeName, &likes.Postid)

		for v := 0; v < len(sdk); v++ {
			if sdk[v].Pid == likes.Postid {
				var count int
				conn.QueryRow("SELECT COUNT(*) FROM likes_ WHERE post_id = ? ", likes.Postid).Scan(&count)

				if count <= 0 {
					sdk[v].Likes = 0
				} else {
					sdk[v].Likes = count
				}
			}
		}
	}

	stD, _ := conn.Query("SELECT user_id, username, post_id FROM dislikes_")
	dislike := reg.Liks{}

	for stD.Next() {

		stD.Scan(&dislike.UserId, &dislike.LikeName, &dislike.Postid)

		for v := 0; v < len(sdk); v++ {
			if sdk[v].Pid == dislike.Postid {
				var count int
				conn.QueryRow("SELECT COUNT(*) FROM dislikes_ WHERE post_id = ? ", dislike.Postid).Scan(&count)

				if count <= 0 {
					sdk[v].Dislikes = 0
				} else if count > 0 {
					sdk[v].Dislikes = count
				}
			}
		}
	}

	stN, eer := conn.Query("SELECT user_id, username, post_id, comment, STRFTIME('%H:%M', date_created) FROM comments")
	if eer != nil {
		log.Fatal(eer)
	}

	var ccc reg.Comments

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

	var newSdk reg.SDK
	newSdk.Us.ID = hydee
	newSdk.Us.Name = usingAname
	newSdk.Us.Education = edu
	newSdk.Us.FullName = fname
	newSdk.Us.HomeAddress = homeAdd
	newSdk.Us.Phone = phoone

	for i := len(sdk) - 1; i >= 0; i-- {
		newSdk.Pst = append(newSdk.Pst, sdk[i])
	}

	fmt.Println("ROWS HERE ->", sdk)
	fmt.Println(w, "Form data saved successfully!")
	reg.Maketmpl(w, "liked", newSdk)

	//Maketmpl(w, "indexlog", sdk)
}

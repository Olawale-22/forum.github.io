package cat

import (
	reg "clouds/getreg"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

//var DBEE = "./datab.db"

func LifeHacks(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Del("Content-Security-Policy")
		w.Header().Set("Cache-Control", "no-cache")

		if r.URL.Path != "/life-hacks" {
			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return
		}

		conn, err := sql.Open("sqlite3", DBEE)
		if err != nil {
			fmt.Println("unable to open database home handler")
			//log.Fatal(err.Error())
		}

		defer conn.Close()

		catty := "life-hacks"

		var sdk []reg.Posts
		sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", catty)
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

		//*******************

		defer sth.Close()

		for sth.Next() {

			var str reg.Posts
			//str.Time = time.Now()
			err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
			if err != nil {
				fmt.Println(err)
				return
			}
			sdk = append(sdk, str)
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
		reg.Maketmpl(w, "life-hacks", newSdk)

		//Maketmpl(w, "indexlog", sdk)
	} else if r.Method == http.MethodPost {

		w.Header().Set("Cache-Control", "no-cache")
		//w.Write([]byte("POST request processed"))
		conn, err := sql.Open("sqlite3", DBEE)
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

		// _, errK := conn.Exec("CREATE TABLE IF NOT EXISTS comments (user_id INTEGER, username TEXT, comment_id INTEGER PRIMARY KEY AUTOINCREMENT, post_id INTEGER, comment TEXT, date_created DATETIME, FOREIGN KEY(user_id) REFERENCES users(user_id), FOREIGN KEY(post_id) REFERENCES activities(post_id))")

		// if errK != nil {
		// 	fmt.Println(errK)
		// }

		var sdk []reg.Posts

		posts := r.FormValue("postit")

		Chat := r.FormValue("chat")

		PostId := r.FormValue("post_IId")

		Like := r.FormValue("lyks")

		Dislike := r.FormValue("unlyks")

		cate := "life-hacks"

		if posts != "" && Chat == "" {
			/****
			_, err = conn.Exec("UPDATE activities SET user_id=?, post=?, username=?, date_created=?", hydee, posts, usingAname, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			****/

			_, err = conn.Exec("INSERT INTO activities (user_id,  post, username, category, date_created) VALUES(?, ?, ?, ?, datetime('now'))", hydee, posts, usingAname, cate, time.Now().Add(time.Hour))
			if err != nil {
				log.Fatal(err)
			}

			sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

			if errNext != nil {
				fmt.Println(errNext)
				return
			}

			defer sth.Close()

			for sth.Next() {

				var str reg.Posts

				//str.Time = time.Now()
				//str.Commenters = Allcmt
				err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
				if err != nil {
					fmt.Println(err)
					return
				}
				sdk = append(sdk, str)
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

		} else if Chat != "" && posts == "" {
			//scrollPosition := r.FormValue("scrollPosition")

			_, err = conn.Exec("INSERT INTO comments (user_id,  username, post_id, comment, date_created) VALUES(?, ?, ?, ?, datetime('now'))", hydee, usingAname, reg.Atoi(PostId), Chat, time.Now().Add(time.Hour))
			if err != nil {
				log.Fatal(err)
			}

			sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities where  category = ?", cate)

			if errNext != nil {
				fmt.Println(errNext)
				return
			}

			defer sth.Close()

			for sth.Next() {

				var str reg.Posts

				//str.Time = time.Now()
				//str.Commenters = Allcmt

				err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
				if err != nil {
					fmt.Println(err)
					return
				}
				sdk = append(sdk, str)
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
			/*var cName, cCmt, cTime string
			var cPid int*/

			// add comments to posts
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
						//ccc.ScrollPosition = Atoi(scrollPosition)
						sdk[v].Commenters = append(sdk[v].Commenters, ccc)
						if len(sdk[v].Commenters) > 8 {
							sdk[v].Commenters = sdk[v].Commenters[0:9]
						}
					}

				}
			}

		} else if Like != "" && Chat == "" && posts == "" && Dislike == "" {

			fmt.Println("#######################--LIIIIIIKE", Like)
			var cty int
			errNN := conn.QueryRow("SELECT COUNT (*) FROM likes_ WHERE user_id = ? AND post_id = ?", hydee, Atoi(Like)).Scan(&cty)

			if errNN != nil {
				fmt.Println(errNN)
			}
			// 	if errNN == sql.ErrNoRows || cty == 0 {
			if cty <= 0 {
				_, err = conn.Exec("INSERT INTO likes_ (user_id, username, post_id) VALUES (?, ?, ?)", hydee, usingAname, Atoi(Like))
				if err != nil {
					log.Fatal(err)
				}

				var count int
				// check if user already like.... to remove like when ever user dislikes
				errNN := conn.QueryRow("SELECT COUNT (*) FROM dislikes_ WHERE user_id = ? AND post_id = ?", hydee, Atoi(Like)).Scan(&count)

				if errNN != nil {
					fmt.Println(errNN)
				}
				if count > 0 {
					_, err = conn.Exec("DELETE FROM dislikes_ WHERE user_id = ? AND username = ? AND post_id= ?", hydee, usingAname, Atoi(Like))
					if err != nil {
						log.Fatal(err)
					}

					sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

					if errNext != nil {
						fmt.Println(errNext)
						return
					}

					defer sth.Close()

					for sth.Next() {

						var str reg.Posts

						//str.Time = time.Now()
						//str.Commenters = Allcmt

						err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
						if err != nil {
							fmt.Println(err)
							return
						}
						sdk = append(sdk, str)
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
								} else if count > 0 {
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

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
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
								//ccc.ScrollPosition = Atoi(scrollPosition)
								sdk[v].Commenters = append(sdk[v].Commenters, ccc)
								if len(sdk[v].Commenters) > 8 {
									sdk[v].Commenters = sdk[v].Commenters[0:9]
								}
							}

						}
					}
				} else {
					// dislike_ count <= 0
					sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

					if errNext != nil {
						fmt.Println(errNext)
						return
					}

					defer sth.Close()

					for sth.Next() {

						var str reg.Posts

						//str.Time = time.Now()
						//str.Commenters = Allcmt

						err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
						if err != nil {
							fmt.Println(err)
							return
						}
						sdk = append(sdk, str)
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
								} else if count > 0 {
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

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
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
								//ccc.ScrollPosition = Atoi(scrollPosition)
								sdk[v].Commenters = append(sdk[v].Commenters, ccc)
								if len(sdk[v].Commenters) > 8 {
									sdk[v].Commenters = sdk[v].Commenters[0:9]
								}
							}

						}
					}
				}

			} else {

				sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

				if errNext != nil {
					fmt.Println(errNext)
					return
				}

				defer sth.Close()

				for sth.Next() {

					var str reg.Posts

					//str.Time = time.Now()
					//str.Commenters = Allcmt

					err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
					if err != nil {
						fmt.Println(err)
						return
					}
					sdk = append(sdk, str)
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
							} else if count > 0 {
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

				/*var cName, cCmt, cTime string
				var cPid int*/

				// add comments to posts
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
							//ccc.ScrollPosition = Atoi(scrollPosition)
							sdk[v].Commenters = append(sdk[v].Commenters, ccc)
							if len(sdk[v].Commenters) > 8 {
								sdk[v].Commenters = sdk[v].Commenters[0:9]
							}
						}

					}
				}
			}
		} else if Dislike != "" && Like == "" && Chat == "" && posts == "" {

			var cty int
			errNN := conn.QueryRow("SELECT COUNT (*) FROM dislikes_ WHERE user_id = ? AND post_id = ?", hydee, Atoi(Dislike)).Scan(&cty)

			if errNN != nil {
				fmt.Println(errNN)
			}

			// 	if errNN == sql.ErrNoRows || cty == 0 {
			if cty <= 0 {
				_, err = conn.Exec("INSERT INTO dislikes_ (user_id, username, post_id) VALUES (?, ?, ?)", hydee, usingAname, Atoi(Dislike))
				if err != nil {
					log.Fatal(err)
				}
				var count int
				// check if user already like.... to remove like when ever user dislikes
				errNN := conn.QueryRow("SELECT COUNT (*) FROM likes_ WHERE user_id = ? AND post_id = ?", hydee, Atoi(Dislike)).Scan(&count)

				if errNN != nil {
					fmt.Println(errNN)
				}
				if count > 0 {
					_, err = conn.Exec("DELETE FROM likes_ WHERE user_id = ? AND username = ? AND post_id= ?", hydee, usingAname, Atoi(Dislike))
					if err != nil {
						log.Fatal(err)
					}
					sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

					if errNext != nil {
						fmt.Println(errNext)
						return
					}

					defer sth.Close()

					for sth.Next() {

						var str reg.Posts

						//str.Time = time.Now()
						//str.Commenters = Allcmt

						err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
						if err != nil {
							fmt.Println(err)
							return
						}
						sdk = append(sdk, str)
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
								} else if count > 0 {
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

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
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
								//ccc.ScrollPosition = Atoi(scrollPosition)
								sdk[v].Commenters = append(sdk[v].Commenters, ccc)
								if len(sdk[v].Commenters) > 8 {
									sdk[v].Commenters = sdk[v].Commenters[0:9]
								}
							}

						}
					}
				} else {

					sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

					if errNext != nil {
						fmt.Println(errNext)
						return
					}

					defer sth.Close()

					for sth.Next() {

						var str reg.Posts

						//str.Time = time.Now()
						//str.Commenters = Allcmt

						err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
						if err != nil {
							fmt.Println(err)
							return
						}
						sdk = append(sdk, str)
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
								} else if count > 0 {
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

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
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
								//ccc.ScrollPosition = Atoi(scrollPosition)
								sdk[v].Commenters = append(sdk[v].Commenters, ccc)
								if len(sdk[v].Commenters) > 8 {
									sdk[v].Commenters = sdk[v].Commenters[0:9]
								}
							}

						}
					}
				}
			} else {

				sth, errNext := conn.Query("SELECT user_id, post_id, post, username, STRFTIME('%H:%M', date_created) FROM activities WHERE  category = ?", cate)

				if errNext != nil {
					fmt.Println(errNext)
					return
				}

				defer sth.Close()

				for sth.Next() {

					var str reg.Posts

					//str.Time = time.Now()
					//str.Commenters = Allcmt

					err = sth.Scan(&str.P_Uid, &str.Pid, &str.Posted, &str.Name, &str.Time)
					if err != nil {
						fmt.Println(err)
						return
					}
					sdk = append(sdk, str)
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
							} else if count > 0 {
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

				/*var cName, cCmt, cTime string
				var cPid int*/

				// add comments to posts
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
							//ccc.ScrollPosition = Atoi(scrollPosition)
							sdk[v].Commenters = append(sdk[v].Commenters, ccc)
							if len(sdk[v].Commenters) > 8 {
								sdk[v].Commenters = sdk[v].Commenters[0:9]
							}
						}

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
		fmt.Fprintf(w, "<html><head><script>window.scrollTo(0, %s);</script></head></html>", PostId)
		reg.Maketmpl(w, "life-hacks", newSdk)
		return
	}
}

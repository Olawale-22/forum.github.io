package getreg

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		ID := r.FormValue("id")

		iD := Atoi(ID)
		conn, err := sql.Open("sqlite3", Deebee)
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

		stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")
		like := Likss{}

		for stLike.Next() {

			stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

			for v := 0; v < len(str.Commenters); v++ {
				if str.Commenters[v].Pid == like.CommentId {
					var count int
					conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

					if count <= 0 {
						str.Commenters[v].Likes = 0
					} else if count > 0 {
						str.Commenters[v].Likes = count
					}
				}
			}
		}
		stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
		dlike := Likss{}

		for stD.Next() {

			stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

			for v := 0; v < len(str.Commenters); v++ {
				if str.Commenters[v].Pid == dlike.CommentId {
					var count int
					conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", dlike.CommentId).Scan(&count)

					if count <= 0 {
						str.Commenters[v].Dislikes = 0
					} else if count > 0 {
						str.Commenters[v].Dislikes = count
					}
				}
			}
		}

		fmt.Println("*************** POST PAGE HERE ->", str)
		fmt.Println(w, "Form data saved successfully!")
		Maketmpl(w, "post", str)
		return

	} else if r.Method == http.MethodPost {

		w.Header().Set("Cache-Control", "no-cache")

		ID := r.FormValue("id")
		iD := Atoi(ID)

		if _, err := strconv.Atoi(ID); err != nil {

			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return
		}

		if r.URL.Path != "/post" || iD < 0 {

			http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
			return

		}

		conn, err := sql.Open("sqlite3", Deebee)
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

		chats := r.FormValue("chats")
		Clikes := r.FormValue("clyks")

		CDislikes := r.FormValue("cdislyks")

		if chats != "" && Clikes == "" && CDislikes == "" {

			_, err = conn.Exec("INSERT INTO comments (user_id,  username, post_id, comment, date_created) VALUES(?, ?, ?, ?, datetime('now'))", str.P_Uid, &str.Name, iD, chats, time.Now().Add(time.Hour))

			if err != nil {
				log.Fatal(err)
			}

			stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

			if eer != nil {
				log.Fatal(eer)
			}

			var ccc Comments
			defer stN.Close()

			for stN.Next() {

				eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

				if eer != nil {
					fmt.Println(eer)
					return
				}
				str.Commenters = append(str.Commenters, ccc)
			}
			stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")
			like := Likss{}

			for stLike.Next() {

				stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

				for v := 0; v < len(str.Commenters); v++ {
					if str.Commenters[v].Cid == like.CommentId {
						var count int
						conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

						if count <= 0 {
							str.Commenters[v].Likes = 0
						} else if count > 0 {
							str.Commenters[v].Likes = count
						}
					}
				}
			}
			stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
			dlike := Likss{}

			for stD.Next() {

				stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

				for v := 0; v < len(str.Commenters); v++ {
					if str.Commenters[v].Cid == dlike.CommentId {
						var count int
						conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", dlike.CommentId).Scan(&count)

						if count <= 0 {
							str.Commenters[v].Dislikes = 0
						} else if count > 0 {
							str.Commenters[v].Dislikes = count
						}
					}
				}
			}

		} else if Clikes != "" && CDislikes == "" && chats == "" {

			var cty int
			errNN := conn.QueryRow("SELECT COUNT (*) FROM clikes_ WHERE user_id = ? AND comment_id = ?", str.P_Uid, Atoi(Clikes)).Scan(&cty)

			if errNN != nil {
				fmt.Println("problem counting from cdlikes clikes")
				fmt.Println(errNN)
			}
			// 	if errNN == sql.ErrNoRows || cty == 0 {
			if cty <= 0 {
				_, err = conn.Exec("INSERT INTO clikes_ (user_id, username, comment_id) VALUES (?, ?, ?)", str.P_Uid, str.Name, Atoi(Clikes))
				if err != nil {
					log.Fatal(err)
				}

				var count int
				// check if user already like.... to remove like when ever user dislikes
				errNN := conn.QueryRow("SELECT COUNT (*) FROM cdislikes_ WHERE user_id = ? AND comment_id = ?", str.P_Uid, Atoi(Clikes)).Scan(&count)

				if errNN != nil {
					fmt.Println(errNN)
				}

				if count > 0 {
					_, err = conn.Exec("DELETE FROM cdislikes_ WHERE user_id = ? AND username = ? AND comment_id= ?", str.P_Uid, str.Name, Atoi(Clikes))
					if err != nil {
						log.Fatal(err)
					}

					stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

					if eer != nil {
						log.Fatal(eer)
					}

					var ccc Comments
					defer stN.Close()

					for stN.Next() {

						eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

						if eer != nil {
							fmt.Println(eer)
							return
						}
						str.Commenters = append(str.Commenters, ccc)
					}

					stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")
					like := Likss{}

					for stLike.Next() {

						stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == like.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Likes = 0
								} else if count > 0 {
									str.Commenters[v].Likes = count
								}
							}
						}
					}

					stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
					dlike := Likss{}

					for stD.Next() {

						stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == dlike.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", dlike.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Dislikes = 0
								} else if count > 0 {
									str.Commenters[v].Dislikes = count
								}
							}
						}
					}

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
				} else {
					// dislike_ count <= 0
					stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

					if eer != nil {
						log.Fatal(eer)
					}

					var ccc Comments
					defer stN.Close()

					for stN.Next() {

						eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

						if eer != nil {
							fmt.Println(eer)
							return
						}
						str.Commenters = append(str.Commenters, ccc)
					}

					stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")

					like := Likss{}

					for stLike.Next() {

						stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == like.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Likes = 0
								} else if count > 0 {
									str.Commenters[v].Likes = count
								}
							}
						}
					}
					stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
					dlike := Likss{}

					for stD.Next() {

						stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == dlike.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", dlike.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Dislikes = 0
								} else if count > 0 {
									str.Commenters[v].Dislikes = count
								}
							}
						}
					}
				}

				/*var cName, cCmt, cTime string
				var cPid int*/

				// add comments to posts

			} else {

				stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

				if eer != nil {
					log.Fatal(eer)
				}

				var ccc Comments
				defer stN.Close()

				for stN.Next() {

					eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

					if eer != nil {
						fmt.Println(eer)
						return
					}
					str.Commenters = append(str.Commenters, ccc)
				}

				stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")

				like := Likss{}

				for stLike.Next() {

					stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

					for v := 0; v < len(str.Commenters); v++ {
						if str.Commenters[v].Cid == like.CommentId {
							var count int
							conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

							if count <= 0 {
								str.Commenters[v].Likes = 0
							} else if count > 0 {
								str.Commenters[v].Likes = count
							}
						}
					}
				}
				stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
				dlike := Likss{}

				for stD.Next() {

					stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

					for v := 0; v < len(str.Commenters); v++ {
						if str.Commenters[v].Cid == dlike.CommentId {
							var count int
							conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", dlike.CommentId).Scan(&count)

							if count <= 0 {
								str.Commenters[v].Dislikes = 0
							} else if count > 0 {
								str.Commenters[v].Dislikes = count
							}
						}
					}
				}
			}

			/*var cName, cCmt, cTime string
			var cPid int*/

			// add comments to posts

		} else if Clikes == "" && CDislikes != "" && chats == "" {

			var cty int
			errNN := conn.QueryRow("SELECT COUNT (*) FROM cdislikes_ WHERE user_id = ? AND comment_id = ?", str.P_Uid, Atoi(CDislikes)).Scan(&cty)

			if errNN != nil {
				fmt.Println("problem counting from clikes cdislikes")
				fmt.Println(errNN)
			}

			// 	if errNN == sql.ErrNoRows || cty == 0 {
			if cty <= 0 {
				_, err = conn.Exec("INSERT INTO cdislikes_ (user_id, username, comment_id) VALUES (?, ?, ?)", str.P_Uid, str.Name, Atoi(CDislikes))
				if err != nil {
					log.Fatal(err)
				}
				var count int
				// check if user already like.... to remove like when ever user dislikes
				errNN := conn.QueryRow("SELECT COUNT (*) FROM clikes_ WHERE user_id = ? AND comment_id = ?", str.P_Uid, Atoi(CDislikes)).Scan(&count)

				if errNN != nil {
					fmt.Println("problem counting from cdislikes")
					fmt.Println(errNN)
				}

				if count > 0 {
					_, err = conn.Exec("DELETE FROM clikes_ WHERE user_id = ? AND username = ? AND comment_id= ?", str.P_Uid, str.Name, Atoi(CDislikes))

					if err != nil {
						log.Fatal(err)
					}

					stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

					if eer != nil {
						log.Fatal(eer)
					}

					var ccc Comments
					defer stN.Close()

					for stN.Next() {

						eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

						if eer != nil {
							fmt.Println(eer)
							return
						}
						str.Commenters = append(str.Commenters, ccc)
					}

					stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")

					like := Likss{}

					for stLike.Next() {

						stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == like.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Likes = 0
								} else if count > 0 {
									str.Commenters[v].Likes = count
								}
							}
						}
					}

					stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
					dlike := Likss{}

					for stD.Next() {

						stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == dlike.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Dislikes = 0
								} else if count > 0 {
									str.Commenters[v].Dislikes = count
								}
							}
						}
					}

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts
				} else {

					stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

					if eer != nil {
						log.Fatal(eer)
					}

					var ccc Comments
					defer stN.Close()

					for stN.Next() {

						eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

						if eer != nil {
							fmt.Println(eer)
							return
						}
						str.Commenters = append(str.Commenters, ccc)
					}

					stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")

					like := Likss{}

					for stLike.Next() {

						stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == like.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Likes = 0
								} else if count > 0 {
									str.Commenters[v].Likes = count
								}
							}
						}
					}

					stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
					dlike := Likss{}

					for stD.Next() {

						stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

						for v := 0; v < len(str.Commenters); v++ {
							if str.Commenters[v].Cid == dlike.CommentId {
								var count int
								conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

								if count <= 0 {
									str.Commenters[v].Dislikes = 0
								} else if count > 0 {
									str.Commenters[v].Dislikes = count
								}
							}
						}
					}

					/*var cName, cCmt, cTime string
					var cPid int*/

					// add comments to posts

				}
			} else {

				stN, eer := conn.Query("SELECT user_id, username, comment_id, comment, STRFTIME('%H:%M', date_created) FROM comments WHERE post_id= ?", iD)

				if eer != nil {
					log.Fatal(eer)
				}

				var ccc Comments
				defer stN.Close()

				for stN.Next() {

					eer = stN.Scan(&ccc.Id, &ccc.Name, &ccc.Cid, &ccc.Comment, &ccc.Time)

					if eer != nil {
						fmt.Println(eer)
						return
					}
					str.Commenters = append(str.Commenters, ccc)
				}

				stLike, _ := conn.Query("SELECT user_id, username, comment_id FROM clikes_")

				like := Likss{}

				for stLike.Next() {

					stLike.Scan(&like.UserId, &like.LikeName, &like.CommentId)

					for v := 0; v < len(str.Commenters); v++ {
						if str.Commenters[v].Cid == like.CommentId {
							var count int
							conn.QueryRow("SELECT COUNT(*) FROM clikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

							if count <= 0 {
								str.Commenters[v].Likes = 0
							} else if count > 0 {
								str.Commenters[v].Likes = count
							}
						}
					}
				}

				stD, _ := conn.Query("SELECT user_id, username, comment_id FROM cdislikes_")
				dlike := Likss{}

				for stD.Next() {

					stD.Scan(&dlike.UserId, &dlike.LikeName, &dlike.CommentId)

					for v := 0; v < len(str.Commenters); v++ {
						if str.Commenters[v].Cid == dlike.CommentId {
							var count int
							conn.QueryRow("SELECT COUNT(*) FROM cdislikes_ WHERE comment_id = ? ", like.CommentId).Scan(&count)

							if count <= 0 {
								str.Commenters[v].Dislikes = 0
							} else if count > 0 {
								str.Commenters[v].Dislikes = count
							}
						}
					}
				}

				/*var cName, cCmt, cTime string
				var cPid int*/

				// add comments to posts

			}
		} else {
			Clikes = ""
			CDislikes = ""
			chats = ""

		}

		fmt.Println("*************** POST PAGE HERE ->", str)
		fmt.Println(w, "Form data saved successfully!")
		Maketmpl(w, "post", str)
		return
	}
}

type Likss struct {
	LikeName  string
	UserId    int
	CommentId int
}

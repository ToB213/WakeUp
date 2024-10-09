package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var tmpl = template.Must(template.ParseFiles("./templates/alarm.html"))

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./wake_up.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", setAlarmHandler)

	http.HandleFunc("/wake_up", wakeUpHandler)

	fmt.Println("サーバーを開始します: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func setAlarmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		alarmTime := r.FormValue("alarm_time")

		_, err := db.Exec("INSERT INTO alarms (time, status) VALUES (?, 0)", alarmTime)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			now := time.Now()
			alarm, _ := time.Parse("15:04", alarmTime) 
			alarmToday := time.Date(now.Year(), now.Month(), now.Day(), alarm.Hour(), alarm.Minute(), 0, 0, now.Location())

			time.Sleep(time.Until(alarmToday))

			var status int
			err = db.QueryRow("SELECT status FROM alarms ORDER BY id DESC LIMIT 1").Scan(&status)
			if err != nil {
				log.Fatal(err)
			}

			if status == 0 {
				err = deleteFilesInDirectory("~/src/wp/test/sample.txt")
				if err != nil {
					log.Fatal("ファイル削除中にエラーが発生しました:", err)
				}
				fmt.Println("時間内に起きなかったため、ファイルを削除しました。")
			}
		}()
	}

	tmpl.Execute(w, nil)
}

func wakeUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_, err := db.Exec("UPDATE alarms SET status = 1 WHERE id = (SELECT MAX(id) FROM alarms)")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(w, "起床確認されました！")
	}
}

func deleteFilesInDirectory(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
		fmt.Println("削除されました:", file)
	}
	return nil
}

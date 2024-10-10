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
	InitDB()
	defer db.Close()

	http.HandleFunc("/", InputHandler)
	http.HandleFunc("/wake_up", WakeUpHandler)
	http.HandleFunc("/create_db", createDBHandler)

	fmt.Println("サーバーを開始します: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./wake_up.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS alarms (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alarm_time TEXT NOT NULL,
			dir_path TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 0 -- 0: 未確認, 1: 確認済み, 2: 削除済み
		);
	`)
	if err != nil {
		log.Fatal("alarmsテーブルの作成に失敗しました: ", err)
	}
	fmt.Println("データベースが初期化されました。")
}

func InputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleFormSubmission(w, r)
	} else {
		tmpl.Execute(w, nil)
	}
}

func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	alarmTime := r.FormValue("alarm_time")
	dirPath := r.FormValue("dir_path")

	// データベースにアラームとディレクトリパスを登録
	_, err := db.Exec("INSERT INTO alarms (alarm_time, dir_path, status) VALUES (?, ?, ?)", alarmTime, dirPath, 0)
	if err != nil {
		log.Println("INSERTに失敗しました: ", err)
		http.Error(w, "データベースへの挿入に失敗しました", http.StatusInternalServerError)
		return
	}

	log.Println("アラームを登録しました: ", alarmTime)
	log.Println("削除対象のディレクトリ: ", dirPath)

	// アラームを設定して、指定された時間後にディレクトリ内のファイルを削除
	go setAlarm(alarmTime, dirPath)
	tmpl.Execute(w, nil)
}

func createDBHandler(w http.ResponseWriter, r *http.Request) {
	InitDB()
	tmpl.Execute(w, nil)
	log.Println("データベースが初期化されました。")
}

func setAlarm(alarmTime string, dirPath string) {
	// アラーム時間を解析
	parsedAlarmTime, err := time.Parse("2006-01-02T15:04", alarmTime)
	if err != nil {
		log.Println("時間の解析に失敗しました: ", err)
		return
	}
	log.Printf("アラームをセットしました: %s", alarmTime)

	now := time.Now()
	alarmToday := time.Date(now.Year(), now.Month(), parsedAlarmTime.Day(), parsedAlarmTime.Hour(), parsedAlarmTime.Minute(), 0, 0, time.Local)

	// アラーム時間までスリープ
	if time.Until(alarmToday) > 0 {
		time.Sleep(time.Until(alarmToday))
	} else {
		log.Println("アラーム時間が既に過ぎています。")
		return
	}

	// 現在のアラームのステータスを取得
	var status int
	err = db.QueryRow("SELECT status FROM alarms WHERE dir_path = ?", dirPath).Scan(&status)
	if err != nil {
		log.Printf("アラームのステータス取得に失敗しました: %s", err)
		return
	}

	if status == 1 || status == 2 {
		log.Printf("アラームは既に確認済みまたは削除済みです: %d", status)
		return
	}	

	// アラーム時間が過ぎたらディレクトリ内のファイルを削除
	log.Printf("アラーム時間を過ぎています: %s", alarmTime)
	err = deleteFilesInDirectory(dirPath)
	if err != nil {
		log.Printf("ディレクトリ内のファイルの削除に失敗しました: %s", err)
	} else {
		log.Printf("ディレクトリ内のファイルを削除しました: %s", dirPath)

		// アラームのステータスを削除済みに更新
		_, err = db.Exec("UPDATE alarms SET status = 2 WHERE dir_path = ?", dirPath)
		if err != nil {
			log.Printf("ステータスの更新に失敗しました: %s", err)
		}
	}

	deleteDB()
}

// deleteFilesInDirectoryは、指定されたディレクトリ内のすべてのファイルを削除する
func deleteFilesInDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		err := os.Remove(filePath)
		if err != nil {
			log.Printf("ファイル削除に失敗しました: %v", err)
		} else {
			log.Printf("ファイル削除成功: %s", filePath)
		}
	}

	return nil
}

// WakeUpHandlerは、起きた確認を処理するハンドラー
func WakeUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// 最新のアラームのステータスを確認済みに変更する
		_, err := db.Exec("UPDATE alarms SET status = 1 WHERE status = 0")
		if err != nil {
			log.Println("ステータスの更新に失敗しました: ", err)
			http.Error(w, "ステータスの更新に失敗しました", http.StatusInternalServerError)
			return
		}

		log.Println("起きた確認が完了しました")
		tmpl.Execute(w, nil)
	}

	deleteDB()
}

func deleteDB() {
	err := os.Remove("./wake_up.db")
	if err != nil {
		log.Println("データベースの削除しました", err)
	}
}


package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// データベース接続
	db, err := sql.Open("sqlite3", "./wake_up.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブルの作成 
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS alarms (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 0  -- 0: 未確認, 1: 確認済み
		);
	`)
	if err != nil {
		log.Fatal("alarmsテーブルの作成に失敗しました: ", err)
	}

	fmt.Println("データベースが初期化されました。")
}

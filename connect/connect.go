package connect

import (
	"database/sql"
	"fmt"
	"log"
)

// Cnn connect
var Cnn *sql.DB

// DatabaceConnect DBに接続
func DatabaceConnect() {

	// MySqlに接続
	var err error
	Cnn, err = sql.Open("mysql", "root:konkonneet21@/testschema") // 使用しているDBの種類, パスワード, スキーマ名

	fmt.Println("\x1b[33m" + "Connected to mysql." + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")

	//接続でエラーが発生した場合の処理
	if err != nil {
		fmt.Println("connect_error")
		log.Fatal(err)
	}
}

// GetCnn データベース情報Get
func GetCnn() *sql.DB {

	return Cnn
}

// CloseDatabace Databaceとの接続終了処理
func CloseDatabace() {

	Cnn.Close()
}

package loadclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../connect"
	"../registernumber"
)

// ClientDataLoad Load用
type ClientDataLoad struct {
	UserID    string `json:"userID"`
	LoginID   string `json:"loginID"`
	Passward  string `json:"passward"`
	DeleteKey int    `json:"deleteKey"`
}

// LoadUserData ユーザー情報が登録されているかcheck
func LoadUserData(w http.ResponseWriter, r *http.Request) {

	// 最終的なユーザー情報を格納するためのスライス
	var loadInfo []ClientDataLoad
	// loadInfoの中身をリセットに戻すために使用
	var loadInfoWork []ClientDataLoad
	// loadInfoにappendしていく時に使用
	var add ClientDataLoad

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	rows, err := getCnn.Query("SELECT * FROM testschema.clientInfo")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for rows.Next() {

		err := rows.Scan(&add.UserID, &add.LoginID, &add.Passward, &add.DeleteKey)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
		loadInfo = append(loadInfo, add)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(loadInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ここで全ユーザーのログインID、パスワードをClientにおくっている。
	// 暗号化して送らないと危険。
	// Jsonにしてから暗号化するべきか、Jsonにする前に暗号化するべきか。@todo
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	// デバッグ用に全クライアント数を取得
	var record = registernumber.GetClientRecordNumber()

	fmt.Println("\n全登録ユーザー情報")
	for i := 0; i < record; i++ {

		fmt.Printf("\n%d人目の情報\n", i+1)
		fmt.Println("ユーザーID: " + "\x1b[33m" + loadInfo[i].UserID + "\x1b[0m")
		fmt.Println("ログインID: " + "\x1b[33m" + loadInfo[i].LoginID + "\x1b[0m")
		fmt.Println("パスワード: " + "\x1b[33m" + loadInfo[i].Passward + "\x1b[0m")
	}
	fmt.Println("\n----------------------------------------------------")

	// charInfoの中身を空に戻す
	loadInfo = loadInfoWork
}

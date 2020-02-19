package loadclient

import (
	"bytes"
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

// ClientDataLoadData Json変換用
type ClientDataLoadData struct {
	Data ClientDataLoad `json:"data"`
}

// LoadUserData ユーザー情報が登録されているかcheck
// 現状クライアント側に全ユーザー情報を送ってしまっているので、それは良くない
func LoadUserData(w http.ResponseWriter, r *http.Request) {

	// デバッグコンソール用
	LoadAllUserData()

	// #region　クライアントから入力されたログインIDとパスワードを取得する
	getLoadData := new(bytes.Buffer)
	getLoadData.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := getLoadData.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	// のちにsavedata.LoginIDとsavedata.passwardに格納されたものと比較する
	var savedata ClientDataLoadData
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// #endregion

	// 最終的なユーザー情報を格納するためのスライス
	var loadInfo []ClientDataLoad
	// loadInfoの中身をリセットに戻すために使用
	var loadInfoWork []ClientDataLoad
	// loadInfoにappendしていく時に使用
	var add ClientDataLoad

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	rows, err := getCnn.Query("SELECT * FROM testschema.clientInfo WHERE loginID = ? AND passward = ?", savedata.Data.LoginID, savedata.Data.Passward)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
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
	//response, err := json.Marshal(Encrypted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ここで全ユーザーのログインID、パスワードをClientにおくっている。
	// 暗号化して送らないと危険。
	// Jsonにしてから暗号化するべきか、Jsonにする前に暗号化するべきか。@todo
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	if loadInfo != nil {

		fmt.Println("\nユーザー情報")
		fmt.Println("\nユーザーID: " + "\x1b[33m" + loadInfo[0].UserID + "\x1b[0m")
		fmt.Println("ログインID: " + "\x1b[33m" + loadInfo[0].LoginID + "\x1b[0m")
		fmt.Println("パスワード: " + "\x1b[33m" + loadInfo[0].Passward + "\x1b[0m")
		fmt.Println("\n----------------------------------------------------")
	} else {

		fmt.Println("\n\x1b[33m新規ユーザー\x1b[0m")
		fmt.Println("\n----------------------------------------------------")
	}

	// charInfoの中身を空に戻す
	loadInfo = loadInfoWork
}

// LoadAllUserData 実際にクライアントには送らないが、デバッグ用に全ユーザーデータを照会する
func LoadAllUserData() {

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

	loadInfo = loadInfoWork
}

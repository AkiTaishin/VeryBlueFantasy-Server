package saveclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"../connect"
	"../createid"
)

// ClientData Client情報格納
type ClientData struct {
	UserID    string `json:"userID"`
	LoginID   string `json:"loginID"`
	Passward  string `json:"passward"`
	DeleteKey int    `json:"deleteKey"`
}

// SaveClientdata Jsonからの変換用
// ClientでSerializeした情報は構造体情報として獲得する
// 中身を同じにしなくてはerrorになってしまうので獲得したい情報ClientDataの構造体を作成する
type SaveClientdata struct {
	Data ClientData `json:"data"`
}

// SaveUserData ClientのIDとパスワードを保存する
func SaveUserData(w http.ResponseWriter, r *http.Request) {

	// rのボディ情報を格納する
	saveNewData := new(bytes.Buffer)
	saveNewData.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveNewData.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata SaveClientdata
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベース情報Get
	var getCnn = connect.GetCnn()

	var create = createid.CreateID()

	// Clientの編成情報をここで更新する
	// いまログインしているユーザーのアカウントのテーブルに保存@todo
	update, err := getCnn.Prepare("INSERT INTO testschema.clientInfo SET userID = ?, loginID = ?, passward = ?, deleteKey = ?")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	// デバッグ用
	// パスワードがpassの時は管理者アカウントにする
	if savedata.Data.Passward == "pass" {

		update.Exec(create, savedata.Data.LoginID, savedata.Data.Passward, 0)
	} else {

		// SQLに変数を直打ちすることはできないのでここで?にした部分に変数を代入
		update.Exec(create, savedata.Data.LoginID, savedata.Data.Passward, savedata.Data.DeleteKey)
	}

	CreateNewFormation(create)

	fmt.Println("\nユーザー新規登録完了")
	fmt.Println("\n新規ユーザーID: " + "\x1b[33m" + create + "\x1b[0m")
	fmt.Println("新規ログインID: " + "\x1b[33m" + savedata.Data.LoginID + "\x1b[0m")
	fmt.Println("新規パスワード: " + "\x1b[33m" + savedata.Data.Passward + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")
}

// CreateNewFormation 新規登録したユーザーのフォーメーション情報を作成する
func CreateNewFormation(userID string) {

	var getCnn = connect.GetCnn()
	recordCount := 0

	// formationテーブルのレコード数を最初に取得
	count, err := getCnn.Query("SELECT COUNT(*) as count FROM testschema.formation")
	defer count.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	//レコード数を recordCount に当てはめる
	for count.Next() {

		err := count.Scan(&recordCount)

		if err != nil {
			fmt.Println("\r\nerror")

			panic(err.Error())
		}
	}
	//

	for i := 0; i < 3; i++ {
		update, err := getCnn.Prepare("INSERT INTO testschema.formation SET adminNumber = ?, number = ?, id = ?, name = ?, asset = ?, attack = ?, hp = ?, element = ?, userID = ?")
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		// SQLに変数を直打ちすることはできないのでここで?にした部分に変数を代入
		update.Exec(recordCount, i, -1, -1, -1, -1, -1, -1, userID)
		recordCount++
	}

}

// ResetClientDataFormationData データベース初期化
// テーブルの内部結合をして、クライアントデータとフォーメーションデータを両方deleteしたい
// 使用する際は管理者データがしっかり昇順に並んでいることを確認してから使用すること！！！
// 使用する際は管理者データがしっかり昇順に並んでいることを確認してから使用すること！！！
func ResetClientDataFormationData(w http.ResponseWriter, r *http.Request) {

	// データベース情報Get
	var getCnn = connect.GetCnn()

	delete, err := getCnn.Query("DELETE A, B FROM testschema.clientInfo AS A INNER JOIN testschema.formation AS B ON A.userID = B.userID WHERE A.deleteKey = 1")
	defer delete.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	fmt.Println("\nResetComplete...")
	fmt.Println("\n----------------------------------------------------")
}

// ResetClientData データベース初期化
func ResetClientData(w http.ResponseWriter, r *http.Request) {

	// データベース情報Get
	var getCnn = connect.GetCnn()

	rows, err := getCnn.Query("DELETE FROM testschema.clientInfo WHERE deleteKey = 1")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\nResetComplete...")
	fmt.Println("\n----------------------------------------------------")
}

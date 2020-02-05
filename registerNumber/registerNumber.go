package registernumber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../connect"
)

// Number 登録者数Load用
type Number struct {
	Key            int `json:"key"`
	RegisterNumber int `json:"registerNumber"`
}

// SaveNumber 登録者数Save用
type SaveNumber struct {
	Data Number `json:"data"`
}

// RecordCount ClientInfoテーブルのレコード数を格納する
type RecordCount struct {
	RecordNumber int `json:"recordNumber"`
}

// LoadRegisterNumber ユーザー情報が登録されているかcheck
func LoadRegisterNumber(w http.ResponseWriter, r *http.Request) {

	var getNumbers []Number
	var getNumber Number
	var resetNumbers []Number

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	// @todo SECECT COUNT(*) FROM testschema.clientInfo で登録されているレコード数を取得するようにする
	// このtodoができたらsaveする必要がなくなり、よりスマートな構造にできるはず
	rows, err := getCnn.Query("SELECT * FROM testschema.register")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for rows.Next() {

		err := rows.Scan(&getNumber.Key, &getNumber.RegisterNumber)

		if err != nil {
			fmt.Println("\r\nerror")

			panic(err.Error())
		}
		getNumbers = append(getNumbers, getNumber)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(getNumbers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println(string(response))

	getNumbers = resetNumbers
}

// CountClientInfoRecordNumber ClientInfoテーブルのレコード数を取得する
func CountClientInfoRecordNumber(w http.ResponseWriter, r *http.Request) {

	// 追加していくスライス
	var getRecordCounts []RecordCount
	// リセット用のスライス
	var resetRecordCounts []RecordCount
	// RecordCount情報を保存する構造体
	var getRecordCount RecordCount

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	// @todo SECECT COUNT(*) FROM testschema.clientInfo で登録されているレコード数を取得するようにする
	// このtodoができたらsaveする必要がなくなり、よりスマートな構造にできるはず
	rows, err := getCnn.Query("SELECT COUNT(*) as count FROM testschema.clientInfo")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	//レコードを getRecordCount に当てはめる
	for rows.Next() {

		err := rows.Scan(&getRecordCount.RecordNumber)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
		getRecordCounts = append(getRecordCounts, getRecordCount)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(getRecordCounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	haveRecordNumber := strconv.Itoa(getRecordCount.RecordNumber)
	fmt.Println("\n全登録ユーザー数: " + "\x1b[33m" + haveRecordNumber + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")

	getRecordCounts = resetRecordCounts
}

// GetClientRecordNumber 他のパッケージのデバッグメッセージを見やすくするために
func GetClientRecordNumber() int {

	var getCnn = connect.GetCnn()
	var record = 0

	rows, err := getCnn.Query("SELECT COUNT(*) as count FROM testschema.clientInfo")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	//レコードを record に当てはめる
	for rows.Next() {

		err := rows.Scan(&record)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
	}

	return record
}

// SaveRegisterNumber 登録者数を上書きする(＋１)
func SaveRegisterNumber(w http.ResponseWriter, r *http.Request) {

	// rのボディ情報を格納する
	saveNewData := new(bytes.Buffer)
	saveNewData.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveNewData.String()
	fmt.Println("dataString_:", dataString)
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata SaveNumber
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベース情報Get
	var getCnn = connect.GetCnn()

	// Clientの編成情報をここで更新する
	// いまログインしているユーザーのアカウントのテーブルに保存@todo
	update, err := getCnn.Prepare("UPDATE testschema.register SET registerNumber = ? WHERE 'key = 0'")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	savedata.Data.RegisterNumber = savedata.Data.RegisterNumber + 1

	// SQLに変数を直打ちすることはできないのでここで?にした部分に変数を代入
	update.Exec(savedata.Data.RegisterNumber)

	fmt.Println("SaveUserData_:", savedata.Data.RegisterNumber)
}

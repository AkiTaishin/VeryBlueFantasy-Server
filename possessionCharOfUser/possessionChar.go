package possessionchar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../connect"
)

// ユーザーがログインした時に所持キャラクターをデータベースから引っ張ってくる処理

// Possession ユーザーごとの所持キャラクター情報
type Possession struct {
	Number int    `json:"number"`
	UserID string `json:"userID"`
	CharID int    `json:"charID"`
}

// PossessionData UserID受け取り用
type PossessionData struct {
	Data Possession `json:"data"`
}

// CharCount ユーザーの所持キャラクター数
type CharCount struct {
	CharNumber int `json:"charNumber"`
}

// GetUserIDandSendCharIDs UserIDのユーザーが所持しているキャラクターをcheckするためにまずユーザーIDを獲得する
func GetUserIDandSendCharIDs(w http.ResponseWriter, r *http.Request) {

	// #region ここでUserIDを取得する
	saveUserID := new(bytes.Buffer)
	saveUserID.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveUserID.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata PossessionData
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ここまでUserIDの取得
	// #endregion

	// #region ここからは対応するUserIDが所持しているキャラ情報を返す
	// 今度はpossessioncharofuserテーブルからSECECTする
	var getCharIDs []Possession
	var getCharID Possession
	var resetCharID []Possession

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	getPossessionChar, err := getCnn.Query("SELECT charID FROM testschema.possessioncharofuser WHERE userID = ?", savedata.Data.UserID)
	defer getPossessionChar.Close()
	if err != nil {
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for getPossessionChar.Next() {

		err := getPossessionChar.Scan(&getCharID.CharID)
		if err != nil {
			fmt.Println("\r\nerror")

			panic(err.Error())
		}
		getCharIDs = append(getCharIDs, getCharID)
	}

	// 獲得したgetCharIDsをJSONに変換
	response, err := json.Marshal(getCharIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	fmt.Println(string(response))

	// リセット
	getCharIDs = resetCharID

	// #endregion
}

// CountPossessionCharNumber Userがキャラクターを何体持っているのか
func CountPossessionCharNumber(w http.ResponseWriter, r *http.Request) {

	// #region ここでUserIDを取得する
	saveUserID := new(bytes.Buffer)
	saveUserID.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveUserID.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata PossessionData
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// #endregion

	// 追加していくスライス
	var getPossessionCharCounts []CharCount
	// リセット用のスライス
	var resetPossessionCharCounts []CharCount
	// キャラクター情報を保存する構造体
	var getPossessionCharCount CharCount

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	rows, err := getCnn.Query("SELECT COUNT(*) as count FROM testschema.possessioncharofuser WHERE userID = ?", savedata.Data.UserID)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	//レコードを getRecordCount に当てはめる
	for rows.Next() {

		err := rows.Scan(&getPossessionCharCount.CharNumber)

		if err != nil {
			fmt.Println("\r\nerror")

			panic(err.Error())
		}
		getPossessionCharCounts = append(getPossessionCharCounts, getPossessionCharCount)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(getPossessionCharCounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	// デバッグ表示用
	haveCharNumber := strconv.Itoa(getPossessionCharCount.CharNumber)

	fmt.Println("\nユーザーID: " + "\x1b[33m" + savedata.Data.UserID + "\x1b[0m")
	fmt.Println("所持キャラクター総数: " + "\x1b[33m" + haveCharNumber + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")

	getPossessionCharCounts = resetPossessionCharCounts
}

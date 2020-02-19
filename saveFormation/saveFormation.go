package saveformation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../connect"
)

// GetFormationData Clientからの編成情報格納
type GetFormationData struct {
	AdminNumber int    `json:"admin"`
	DataNumber  int    `json:"number"`
	CharID      int    `json:"id"`
	CharName    string `json:"name"`
	Asset       string `json:"asset"`
	Attack      int    `json:"attack"`
	Hp          int    `json:"hp"`
	Element     int    `json:"element"`
	UserID      string `json:"userID"`
}

// Savedata Jsonからの変換用
// ClientでSerializeした情報は構造体情報として獲得する
// 中身を同じにしなくてはerrorになってしまうので獲得したい情報GetFormationDataの構造体を作成する
type Savedata struct {
	Data GetFormationData `json:"data"`
}

// SaveFormation Clientからのリクエストデータを保存する
func SaveFormation(w http.ResponseWriter, r *http.Request) {

	// rのボディ情報を格納する
	saveNewData := new(bytes.Buffer)
	saveNewData.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveNewData.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata Savedata
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベース情報Get
	var getCnn = connect.GetCnn()

	// Clientの編成情報をここで更新する
	// いまログインしているユーザーのアカウントのテーブルに保存@todo
	update, err := getCnn.Prepare("UPDATE testschema.formation SET number = ?, id = ?, name = ?, asset = ?, attack = ?, hp = ?, element = ? WHERE userID = ? and number = ?")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	// SQLに変数を直打ちすることはできないのでここで?にした部分に変数を代入
	update.Exec(savedata.Data.DataNumber, savedata.Data.CharID, savedata.Data.CharName, savedata.Data.Asset, savedata.Data.Attack, savedata.Data.Hp, savedata.Data.Element, savedata.Data.UserID, savedata.Data.DataNumber)

	dataNumber := strconv.Itoa(savedata.Data.DataNumber + 1)
	hp := strconv.Itoa(savedata.Data.Hp)
	attack := strconv.Itoa(savedata.Data.Attack)

	fmt.Println("\n編成したキャラクター情報")

	fmt.Println("\n配置する場所: " + "\x1b[33m" + dataNumber + "番目\x1b[0m")
	fmt.Println("キャラクター名: " + "\x1b[33m" + savedata.Data.CharName + "\x1b[0m")
	fmt.Println("キャラクター体力: " + "\x1b[33m" + hp + "\x1b[0m")
	fmt.Println("キャラクター攻撃力: " + "\x1b[33m" + attack + "\x1b[0m")
	fmt.Println("AssetPass: " + "\x1b[33m" + savedata.Data.Asset + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")

}

// GetFormation 保存されている編成情報をデータベースから取り出す
func GetFormation(w http.ResponseWriter, r *http.Request) {

	// #region ここでUserIDを取得する
	saveUserID := new(bytes.Buffer)
	saveUserID.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveUserID.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// getUseID = 生データに戻した情報を格納する
	var GetUseID Savedata
	err := json.Unmarshal(dataArray, &GetUseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// #endregion

	// 現在保存されている編成情報をひとつずつ格納していくためのスライス
	var savedata []Savedata
	// リセット用
	var savedataWork []Savedata
	// savedataにapeendしていくための変数
	var add Savedata

	// データベース情報Get
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	// userIDと一致するClientの編成を獲得する
	rows, err := getCnn.Query("SELECT * FROM testschema.formation WHERE userID = ?", GetUseID.Data.UserID)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for rows.Next() {

		err := rows.Scan(&add.Data.AdminNumber, &add.Data.DataNumber, &add.Data.CharID, &add.Data.CharName, &add.Data.Asset, &add.Data.Attack, &add.Data.Hp, &add.Data.Element, &add.Data.UserID)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
		savedata = append(savedata, add)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	noData := strconv.Itoa(-1)

	// デバッグ用の表示メッセージ
	fmt.Println("\n編成キャラクター情報")
	if savedata != nil {
		for i := 0; i < 3; i++ {

			// 編成されている場合
			if savedata[i].Data.CharName != noData {

				fmt.Printf("\n%d番目のキャラクター名: "+"\x1b[33m"+savedata[i].Data.CharName+"\x1b[0m", i+1)

			} else {

				// 編成されていない場合
				fmt.Printf("\n%d番目のキャラクター名: "+"\x1b[33mなし\x1b[0m", i+1)

			}
		}
	} else {

		// 新規登録ユーザーの場合savedataの中身はnil
		for i := 0; i < 3; i++ {

			// 編成されていない場合
			fmt.Printf("\n%d番目のキャラクター名: "+"\x1b[33mなし\x1b[0m", i+1)
		}
	}

	fmt.Println("\n\n----------------------------------------------------")

	// 念のためリセット
	savedata = savedataWork
}

// ResetFormation データベース初期化
func ResetFormation(w http.ResponseWriter, r *http.Request) {

	// #region ここでUserIDを取得する
	saveUserID := new(bytes.Buffer)
	saveUserID.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveUserID.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// getUseID = 生データに戻した情報を格納する
	var getUseID Savedata
	err := json.Unmarshal(dataArray, &getUseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ここまでUserIDの取得
	// #endregion

	// データベース情報Get
	var getCnn = connect.GetCnn()

	// すべて初期化
	for i := 0; i < 3; i++ {

		update, err := getCnn.Prepare("UPDATE testschema.formation SET id = ?, name = ?, asset = ?, attack = ?, hp = ?, element = ? WHERE userID = ?")
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		update.Exec(-1, -1, -1, -1, -1, -1, getUseID.Data.UserID)
	}

	fmt.Println("\nResetComplete...")
	fmt.Println("\n----------------------------------------------------")
}

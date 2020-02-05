package getcharinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"../connect"
)

// CharacterInfo キャラ情報取得
type CharacterInfo struct {
	// 結合元のキャラ情報テーブル
	CharID    int    `json:"id"`
	CharName  string `json:"name"`
	Asset     string `json:"asset"`
	Detail    string `json:"detail"`
	Formation string `json:"formation"`
	Attack    int    `json:"attack"`
	Hp        int    `json:"hp"`
	Element   int    `json:"element"`

	// 結合するキャラ所持情報テーブル
	Number          int    `json:"number"`
	UserID          string `json:"userID"`
	HaveCharacterID int    `json:"charID"`
}

// CharIDsData 変換用
type CharIDsData struct {
	Data CharacterInfo `json:"data"`
}

// GetCharInfo ユーザーが所持しているキャラクターだけ送信
func GetCharInfo(w http.ResponseWriter, r *http.Request) {

	//#region ここで所持しているキャラクターのIDを受け取る為に誰が持っている情報なのか獲得する
	saveCharID := new(bytes.Buffer)
	saveCharID.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveCharID.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata CharIDsData
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ここまでUserIDの取得
	//#endregion

	// 最終的な編成情報を格納するためのスライス
	var charInfo []CharacterInfo
	// charInfoの中身をリセットに戻すために使用
	var charInfoWork []CharacterInfo
	// charInfoにappendしていく時に使用
	var add CharacterInfo

	// データベース情報のGet
	var getCnn = connect.GetCnn()

	// characterinfoテーブルとpossessioncharofuserテーブルを内部結合する
	// その後、characterinfoテーブルのidとpossessioncharofuserテーブルのcharIDが等しいもの（ =そのユーザーが所持しているキャラクター ）を取得する
	rows, err := getCnn.Query("SELECT * FROM testschema.characterinfo INNER JOIN testschema.possessioncharofuser ON testschema.characterinfo.id = testschema.possessioncharofuser.charID WHERE testschema.possessioncharofuser.userID = ?", savedata.Data.UserID)
	// rows, err := getCnn.Query("SELECT FROM testschema.characterinfo INNER JOIN testschema.possessioncharofuser ON testschema.characterinfo.id = testschema.possessioncharofuser.charID WHERE testschema.possessioncharofuser.userID = ?", savedata.Data.UserID)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for rows.Next() {

		err := rows.Scan(&add.CharID, &add.CharName, &add.Asset, &add.Detail, &add.Formation, &add.Attack, &add.Hp, &add.Element, &add.Number, &add.UserID, &add.CharID)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
		charInfo = append(charInfo, add)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(charInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	// charInfoの中身を空に戻す
	charInfo = charInfoWork
}

package scene

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"../connect"
)

// Scene シーン遷移用
type Scene struct {
	Number       int    `json:"number"`
	TemplateName string `json:"templateName"`
}

// GetSceneNumber 受信データを生データに変換する
type GetSceneNumber struct {
	Data Scene `json:"data"`
}

// GoToScene いまある画面遷移先を全て送信
func GoToScene(w http.ResponseWriter, r *http.Request) {

	// #region 返したいシーンのnumberを取り出す
	saveSceneNumber := new(bytes.Buffer)
	saveSceneNumber.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	dataString := saveSceneNumber.String()
	dataArray := []byte(dataString)

	// 獲得したJsonデータを認識できるように生データに変換しなおす
	// dataArray = Jsonデータが現在入っている
	// savedata = 生データに戻した情報を格納する
	var savedata GetSceneNumber
	err := json.Unmarshal(dataArray, &savedata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ここまでシーンnumberの取得
	// #endregion

	// データベースのScene情報をひとつずつ格納していくためのスライス
	var scene []Scene
	// スライスの中身を空に戻すための変数
	var sceneWork []Scene
	// この変数を用いてsceneにひとつずつ格納（append）していく
	var add Scene

	// データベース情報Get
	var getCnn = connect.GetCnn()

	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	rows, err := getCnn.Query("SELECT * FROM testschema.changeScene WHERE number = ?", savedata.Data.Number)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for rows.Next() {

		err := rows.Scan(&add.Number, &add.TemplateName)

		if err != nil {
			fmt.Println(err)

			panic(err.Error())
		}
		scene = append(scene, add)
	}

	// 構造体をJSONに変換
	response, err := json.Marshal(scene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	// デバッグ用の表示メッセージ
	fmt.Println("\n遷移先シーン名: " + "\x1b[33m" + scene[0].TemplateName + "\x1b[0m")
	fmt.Println("\n----------------------------------------------------")

	// sceneの中身のリセット
	scene = sceneWork
}

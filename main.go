package main

import (
	"fmt"
	"net/http"

	message "./Message"
	"./aes"
	"./connect"
	"./getcharinfo"
	loadclient "./loadClient"
	possessionchar "./possessionCharOfUser"
	"./registernumber"
	saveclient "./saveClient"
	"./saveformation"
	"./scene"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	fmt.Println("\n----------------------------------------------------")
	fmt.Println("\n\x1b[33m" + "Running..." + "\x1b[0m")

	// データベースとの接続
	connect.DatabaceConnect()

	// 自分が設定したURIがたたかれたら
	http.HandleFunc("/pleaseSaveRegisterNumber", registernumber.SaveRegisterNumber)    // registernumber保存用
	http.HandleFunc("/pleaseLoadRegisterNumber", registernumber.LoadRegisterNumber)    // registernumber読み込み用
	http.HandleFunc("/pleaseRecordNumber", registernumber.CountClientInfoRecordNumber) // clientInfoテーブルのレコード数取得用

	http.HandleFunc("/pleaseSaveClient", saveclient.SaveUserData)                              // UserData保存用
	http.HandleFunc("/pleaseLoadClientInfo", loadclient.LoadUserData)                          // UserData読み出し
	http.HandleFunc("/pleaseResetClient", saveclient.ResetClientData)                          // UserDataリセット
	http.HandleFunc("/pleaseResetClientAndFormation", saveclient.ResetClientDataFormationData) // ClientテーブルとFormationテーブルリセット

	http.HandleFunc("/pleasePossessionNumber", possessionchar.CountPossessionCharNumber) // UserIDを獲得し、そのユーザーのキャラクター所持数を送る
	http.HandleFunc("/pleaseSendUserID", possessionchar.GetUserIDandSendCharIDs)         // UserIDを獲得し、そのユーザーのキャラクター所持状況を送る

	http.HandleFunc("/goToScene", scene.GoToScene) // 画面遷移用

	http.HandleFunc("/pleaseCharacterInfo", getcharinfo.GetCharInfo)  // キャラクター情報取得用
	http.HandleFunc("/pleaseSaveDetail", saveformation.SaveFormation) // 編成の保存

	http.HandleFunc("/pleaseFormationData", saveformation.GetFormation) // 保存されている編成の読み出し
	http.HandleFunc("/pleaseReset", saveformation.ResetFormation)       // 編成のリセット

	http.HandleFunc("/pleaseResponse", message.SendMessage) // Debug用メッセージの取得
	http.HandleFunc("/pleaseCipher", aes.CipherMessage)     // 暗号/復号
	http.HandleFunc("/pleaseCBC", aes.CBCMessage)           // 暗号/復号（CBC）
	http.HandleFunc("/pleaseCTR", aes.CTRMessage)           // 暗号/復号（CTR）

	// 応答待ち
	http.ListenAndServe(":8080", nil)

	// データベースとの接続の終了処理
	connect.CloseDatabace()
}

package message

import (
	"fmt"
	"net/http"
)

// SendMessage メッセージを送りたい
func SendMessage(w http.ResponseWriter, r *http.Request) {

	// #region コメントアウト
	// 送りたいメッセージ
	data := "{\"test\":\"SuccessSendMessage\"}"
	//data := "{test:SuccessSendMessage}"
	// 書き込むためにバイトに変換
	responses := []byte(data)

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(responses)
	fmt.Println(string(responses))

	//log.Print("[送信完了 : ", w, "]")
	// #endregion

}

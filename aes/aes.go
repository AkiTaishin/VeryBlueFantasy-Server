package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
)

// CipherMessage 送りたいメッセージをAESを用いて暗号化する
func CipherMessage(w http.ResponseWriter, r *http.Request) {

	// 鍵の長さは 16, 24, 32 バイトのどれかにしないとエラー
	// これが共通鍵になる
	key := []byte("aes-secret-key-1")
	// cipher.Block を実装している AES 暗号化オブジェクトを生成する
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// #region 暗号するものについて色々

	// ここに格納するものを暗号化する
	// 暗号化される平文の長さは 16 バイト (128 ビット)
	// 16バイトないとエラーになる
	// 16バイトを越えると16バイト分のみ暗号化する

	// plainText := []byte("髙橋千顕 GN31") ⇒出力　髙橋千顕 GN3
	// 平文の理由16バイトオーバー

	// plainText := []byte("髙橋千顕 GN")	⇒出力　めっちゃエラー
	// 平文が16バイト未満

	// plainText := []byte("髙橋千顕GN31")	⇒出力　OK
	// 平文16バイト丁度

	// つまり、送りたいメッセージは基本的に16バイトを超えると想定されるのでこの方法を拡張していかなければならない
	// ⇒CBC (Cipher Block Chaining)

	// #endregion

	plainText := []byte("髙橋千顕GN31")
	// 暗号化されたバイト列を格納するスライスを用意する
	encrypted := make([]byte, aes.BlockSize)
	// AES で暗号化をおこなう
	c.Encrypt(encrypted, plainText)
	// 結果は暗号化されている
	fmt.Println(string(encrypted))
	// Output:
	// #^ϗ~:f9��˱�1�

	// 復号する
	decrypted := make([]byte, aes.BlockSize)
	c.Decrypt(decrypted, encrypted)
	// 結果は元の平文が得られる
	fmt.Println(string(decrypted))
	// Output:
	// secret plain txt
}

// CBCMessage 上の関数の平文16バイト以上の時の対応策その１
func CBCMessage(w http.ResponseWriter, r *http.Request) {

	// #region CBCについてもあれこれ

	// "secret text 9999"
	// これは16バイト

	// plainText := []byte("髙橋千顕 GN31")
	// パディングしていないのでたくさんエラーが出た

	// #endregion

	// 平文。長さが 16 バイトの整数倍でない場合はパディングする必要がある
	plainText := []byte("secret text 9999")
	fmt.Printf("%d\n", len(plainText))

	// 暗号化データ。先頭に初期化ベクトル (IV) を入れるため、1ブロック分余計に確保する
	encrypted := make([]byte, aes.BlockSize+len(plainText))

	// IV は暗号文の先頭に入れておくことが多い
	iv := encrypted[:aes.BlockSize]
	// IV としてランダムなビット列を生成する
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	// ブロック暗号として AES を使う場合
	key := []byte("secret-key-12345")
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// CBC モードで暗号化する
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], plainText)
	fmt.Printf("encrypted: %x\n", encrypted)

	// 復号するには復号化用オブジェクトが別に必要
	mode = cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(plainText))
	// 先頭の IV を除いた部分を復号する
	mode.CryptBlocks(decrypted, encrypted[aes.BlockSize:])
	fmt.Printf("decrypted: %s\n", decrypted)
	// Output:
	// decrypted: secret text 9999
}

// CTRMessage CTRを用いて暗号化する
// これは平文が16バイトである制限はない
func CTRMessage(w http.ResponseWriter, r *http.Request) {

	secletText := "頼むから暗号化して"
	fmt.Printf("\n平文のバイト数: %d\n\n", len(secletText))

	// 平文 16バイト制限なし
	// つまりパディングの必要はない！！！
	plainText := []byte("頼むから暗号化して")

	// 共通鍵
	// key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	// 鍵にするものに空白は許されない
	key := []byte("sendCipherPleaseDecipher")

	// Create new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// Create IV
	// 1つ前の暗号文ブロックと平文ブロックの内容を混ぜ合わせてから暗号化をおこなう
	// 先頭のブロックはその前のブロックが存在していないのでここで仮想のブロックを作成する
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// Encrypt ここで暗号化
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	fmt.Printf("Cipher text: %x \n", cipherText)

	// Decrpt ここで復元
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	fmt.Printf("Decrypted text: %s\n", string(decryptedText))
}

// // TestCTRMessage テスト
// func TestCTRMessage(getCipherText loadclient.ClientDataLoad) {

// 	// 平文 16バイト制限なし
// 	// つまりパディングの必要はない！！！
// 	seclet := getCipherText.UserID + getCipherText.LoginID + getCipherText.Passward + strconv.Itoa(getCipherText.DeleteKey)
// 	plainText := []byte(seclet)

// 	// 共通鍵
// 	// key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

// 	// 鍵にするものに空白は許されない
// 	key := []byte("sendCipherPleaseDecipher")

// 	// Create new AES cipher block
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Printf("err: %s\n", err)
// 	}

// 	// Create IV
// 	// 1つ前の暗号文ブロックと平文ブロックの内容を混ぜ合わせてから暗号化をおこなう
// 	// 先頭のブロックはその前のブロックが存在していないのでここで仮想のブロックを作成する
// 	cipherText := make([]byte, aes.BlockSize+len(plainText))
// 	iv := cipherText[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		fmt.Printf("err: %s\n", err)
// 	}

// 	// Encrypt ここで暗号化
// 	encryptStream := cipher.NewCTR(block, iv)
// 	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
// 	fmt.Printf("Cipher text: %x \n", cipherText)

// 	// Decrpt ここで復元
// 	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
// 	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
// 	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
// 	fmt.Printf("Decrypted text: %s\n", string(decryptedText))
// }

package createid

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// UserInfo ユーザー情報
type UserInfo struct {
	Name      string
	ID        string
	HighScore int
}

// UserInfoArray ↑の配列
type UserInfoArray []*UserInfo

var userInfoArray UserInfoArray

const (
	rs5Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs5LetterIdxBits = 6
	rs5LetterIdxMask = 1<<rs5LetterIdxBits - 1
	rs5LetterIdxMax  = 63 / rs5LetterIdxBits
)

// CreateID ランダムにIDを振り分け
func CreateID() string {

	runes := make([]byte, 6)

	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(255))
		runes[i] = byte(num.Int64())
	}

	return base64.RawStdEncoding.EncodeToString(runes)

}

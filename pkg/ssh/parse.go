package ssh

import (
	"fmt"

	"github.com/gliderlabs/ssh"
)

func ParsePubKeyString(pubKeyString string) (ssh.PublicKey, error) {
	// SSH公開鍵文字列をバイト列にデコード
	// バイト列をssh.PublicKeyに変換
	decodedKey := []byte(pubKeyString)
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(decodedKey)
	if err != nil {
		return nil, fmt.Errorf("Error parsing SSH public key: %v", err)
	}
	return pubKey, nil
}

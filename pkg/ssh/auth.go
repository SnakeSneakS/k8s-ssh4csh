package ssh

import (
	"log"

	"github.com/gliderlabs/ssh"
)

func PublicKeyHandler(publicKey ssh.PublicKey) ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		if !ssh.KeysEqual(key, publicKey) {
			log.Println("public key don't match")
			return false
		}

		return true
	}
}

func PasswordHandler() ssh.PasswordHandler {
	return func(ctx ssh.Context, password string) bool {
		log.Print("We don't want password")

		return false
	}
}

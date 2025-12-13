package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/curve25519"
)

func GenerateKeyPair() (string, string, error) {
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)

	if err != nil {
		return "", "", err
	}

	//WireGuard uses Curve25519
	//Applying clamping to private key (Wireguard requirement)
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64

	var publicKey [32]byte

	curve25519.ScalarBaseMult(&publicKey, (*[32]byte)(privateKey))

	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey[:])

	return privateKeyBase64, publicKeyBase64, nil
}

func ValidateKeyPair(privateKeyBase64, publicKeyBase64 string) (bool, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return false, err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, err
	}

	if len(privateKey) != 32 || len(publicKey) != 32 {
		return false, errors.New("invalid key length")
	}

	var generatedPublicKey [32]byte

	curve25519.ScalarBaseMult(&generatedPublicKey, (*[32]byte)(privateKey))

	for i := 0; i < 32; i++ {
		if generatedPublicKey[i] != publicKey[i] {
			return false, nil
		}
	}
	return true, nil
}

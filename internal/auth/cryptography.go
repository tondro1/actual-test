package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(str string) ([]byte, error) {
	byteStr := []byte(str)
	key, err := getSecretKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key.N.Bytes())
	if err != nil {
		return nil, err
	}

	b := base64.StdEncoding.EncodeToString(byteStr)
	
	cipherText := make([]byte, aes.BlockSize + len(b))

	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(b))

	return cipherText, nil
}
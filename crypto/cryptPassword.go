package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/Archie1978/regate/configuration"
)

var (
	initialVector = "1010101010101010"
	passphrase    = []byte{0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31}
)

func AESEncrypt(src string, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return crypted
}
func AESDecrypt(crypt []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return PKCS5Trimming(decrypted)
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func DecryptPasswordString(data string) string {
	if strings.HasPrefix(data, "aes:") && len(configuration.ConfigurationGlobal.KeyCrypt) > 0 {
		base64Code, err := base64.StdEncoding.DecodeString(data[4:])
		if err != nil {
			return ""
		}
		return string(AESDecrypt(base64Code, configuration.ConfigurationGlobal.KeyCrypt))
	}
	return data
}

func CryptPasswordString(data string) string {
	if len(configuration.ConfigurationGlobal.KeyCrypt) > 0 {
		log.Fatal("KeyCrypt not valid")
	}
	return "aes:" + base64.StdEncoding.EncodeToString(AESEncrypt(data, configuration.ConfigurationGlobal.KeyCrypt))
}

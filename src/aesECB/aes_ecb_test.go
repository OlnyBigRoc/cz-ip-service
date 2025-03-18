package aesECB

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	key := "1234567890123456" // AES-128 需要16字节密钥
	src := "Hello, World!"

	encrypted := EncryptByAes(src, key)
	fmt.Println("Encrypted:", encrypted)

	decrypted := DecryptByAes(encrypted, key)
	fmt.Println("Decrypted:", decrypted)
}

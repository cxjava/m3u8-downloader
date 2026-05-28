package decrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestDecrypt(t *testing.T) {
	key := []byte("1234567890123456") // 16-byte AES-128 key
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	plaintext := []byte("hello world, this is a test message for aes cbc")

	// encrypt manually with AES-CBC + PKCS7 padding
	padded := pkcs7Pad(plaintext, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	ciphertext := make([]byte, len(padded))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, padded)

	// now test decryption
	got, err := Decrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, plaintext) {
		t.Errorf("Decrypt() = %v, want %v", got, plaintext)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key := []byte("1234567890123456")
	wrongKey := []byte("abcdefghijklmnop")
	iv := make([]byte, 16)

	plaintext := []byte("test data")
	padded := pkcs7Pad(plaintext, aes.BlockSize)
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(padded))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, padded)

	result, err := Decrypt(ciphertext, wrongKey, iv)
	if err != nil {
		t.Fatal(err)
	}
	// wrong key produces garbage output (not the original plaintext)
	if bytes.Equal(result, plaintext) {
		t.Error("Decrypt() with wrong key should not return original plaintext")
	}
}

func TestPKCS7UnPadding(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "unpad 1 byte",
			data: []byte("hello world\x01"),
			want: []byte("hello world"),
		},
		{
			name: "unpad 5 bytes",
			data: []byte("hello\x05\x05\x05\x05\x05"),
			want: []byte("hello"),
		},
		{
			name: "unpad full block",
			data: append([]byte("test"), bytes.Repeat([]byte{12}, 12)...),
			want: []byte("test"),
		},
		{
			name: "no padding needed (unpad 16)",
			data: bytes.Repeat([]byte{16}, 16),
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PKCS7UnPadding(tt.data)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PKCS7UnPadding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

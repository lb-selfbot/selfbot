package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type Encryption struct {
	Key string
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func pad(text []byte, blocksize int) []byte {
	// Calculate the padding needed
	padding := blocksize - (len(text) % blocksize)

	// Create the padding
	for i := 0; i < padding; i++ {
		text = append(text, byte(padding))
	}

	return text
}

func unpad(text []byte, blocksize int) []byte {
	// Get the padding
	padding := int(text[len(text)-1])

	// Remove the padding
	text = text[:len(text)-padding]

	return text
}

func (e Encryption) Encrypt(text []byte) ([]byte, error) {
	// Generate a random salt
	salt, err := randomBytes(32)
	if err != nil {
		return nil, err
	}

	// Hash the key with the salt
	key := pbkdf2.Key([]byte(e.Key), salt, 15000, 32, sha256.New)

	// Generate a random IV
	iv, err := randomBytes(16)
	if err != nil {
		return nil, err
	}

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Pad the text
	paddedText := pad(text, block.BlockSize())

	// Create the cipher
	ciphertext := make([]byte, len(paddedText))

	// Encrypt the text
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// Return the encrypted text
	return []byte(fmt.Sprintf(
		"%s:%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(ciphertext),
	)), nil
}

func (e Encryption) Decrypt(text []byte) ([]byte, error) {
	// Split the text into the salt, IV, and ciphertext
	splitText := strings.Split(string(text), ":")
	salt, err := base64.StdEncoding.DecodeString(splitText[0])
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(splitText[1])
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(splitText[2])
	if err != nil {
		return nil, err
	}

	// Hash the key with the salt
	key := pbkdf2.Key([]byte(e.Key), salt, 15000, 32, sha256.New)

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Create the cipher
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the text
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Return the decrypted text
	return unpad(plaintext, block.BlockSize()), nil
}

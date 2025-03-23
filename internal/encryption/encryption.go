package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

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
	for range padding {
		text = append(text, byte(padding))
	}

	return text
}

func unpad(text []byte) []byte {
	// Get the padding
	padding := int(text[len(text)-1])

	// Remove the padding
	text = text[:len(text)-padding]

	return text
}

type AESCBC struct {
	Key string
}

func (a AESCBC) Encrypt(text []byte) ([]byte, error) {
	// Generate a random salt
	salt, err := randomBytes(32)
	if err != nil {
		return nil, err
	}

	// Hash the key with the salt
	key := pbkdf2.Key([]byte(a.Key), salt, 15000, 32, sha256.New)

	// Generate a random IV
	iv, err := randomBytes(16)
	if err != nil {
		return nil, err
	}

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad the text
	paddedText := pad(text, block.BlockSize())

	// Create the cipher
	ciphertext := make([]byte, len(paddedText))

	// Encrypt the text
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// Return the encrypted text
	return fmt.Appendf(nil,
		"%s:%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(ciphertext),
	), nil
}

func (a AESCBC) Decrypt(text []byte) ([]byte, error) {
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
	key := pbkdf2.Key([]byte(a.Key), salt, 15000, 32, sha256.New)

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create the cipher
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the text
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Return the decrypted text
	return unpad(plaintext), nil
}

type AES struct {
	Key string
}

func (a AES) Encrypt(text string, keyLen int) (string, error) {
	// Generate a random salt
	salt, err := randomBytes(16)
	if err != nil {
		return "", err
	}

	// Hash the key with the salt
	key := pbkdf2.Key([]byte(a.Key), salt, 15000, keyLen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce (number used once)
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and seal the data
	cipherText := gcm.Seal(nonce, nonce, []byte(text), nil)

	// Return the salt and encrypted data
	return base64.StdEncoding.EncodeToString(salt) + "." + base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a AES) Decrypt(text string, keyLen int) (string, error) {
	split := strings.Split(text, ".")
	if len(split) != 2 {
		return "", fmt.Errorf("invalid ciphertext format")
	}

	salt, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(a.Key), salt, 15000, keyLen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := gcm.NonceSize()

	// Make sure the ciphertext is long enough
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce from the ciphertext
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt the data
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

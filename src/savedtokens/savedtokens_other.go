//go:build !windows

package savedtokens

func win32CryptUnprotectData(data []byte) ([]byte, error) { return nil, nil }

func searchEncryptedTokens(content []byte, savedTokens *SavedTokens) {}

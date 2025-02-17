//go:build !windows

package savedtokens

func searchEncryptedTokens(content []byte, savedTokens *SavedTokens) {}

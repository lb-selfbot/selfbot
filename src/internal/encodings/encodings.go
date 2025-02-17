package encodings

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"net/url"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/encryption"
	"github.com/darkwyrm/b85"
	"github.com/mr-tron/base58/base58"
	"github.com/multiformats/go-base36"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

type Encoding struct {
	Encode func(string) (string, error)
	Decode func(string) (string, error)
}

var morseEncoding = [][]string{
	{"A", ".-"},
	{"B", "-..."},
	{"C", "-.-."},
	{"D", "-.."},
	{"E", "."},
	{"F", "..-."},
	{"G", "--."},
	{"H", "...."},
	{"I", ".."},
	{"J", ".---"},
	{"K", "-.-"},
	{"L", ".-.."},
	{"M", "--"},
	{"N", "-."},
	{"O", "---"},
	{"P", ".--."},
	{"Q", "--.-"},
	{"R", ".-."},
	{"S", "..."},
	{"T", "-"},
	{"U", "..-"},
	{"V", "...-"},
	{"W", ".--"},
	{"X", "-..-"},
	{"Y", "-.--"},
	{"Z", "--.."},
	{"0", "-----"},
	{"1", ".----"},
	{"2", "..---"},
	{"3", "...--"},
	{"4", "....-"},
	{"5", "....."},
	{"6", "-...."},
	{"7", "--..."},
	{"8", "---.."},
	{"9", "----."},
}

var Encodings = map[string]Encoding{
	"rot13": {
		Encode: func(s string) (string, error) {
			return strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' {
					return 'a' + (r-'a'+13)%26
				}
				if r >= 'A' && r <= 'Z' {
					return 'A' + (r-'A'+13)%26
				}
				return r
			}, s), nil
		},
	},
	"hex": {
		Encode: func(s string) (string, error) {
			return hex.EncodeToString([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := hex.DecodeString(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"binary": {
		Encode: func(s string) (string, error) {
			var buf strings.Builder
			for _, c := range s {
				_, err := buf.WriteString(strings.ReplaceAll(fmt.Sprintf("%08b", c), " ", "0"))
				if err != nil {
					return "", err
				}
			}
			return buf.String(), nil
		},
		Decode: func(s string) (string, error) {
			var buf strings.Builder
			for i := 0; i+8 <= len(s); i += 8 {
				b, err := strconv.ParseUint(s[i:i+8], 2, 8)
				if err != nil {
					return "", err
				}
				buf.WriteByte(byte(b))
			}
			return buf.String(), nil
		},
	},
	"base32": {
		Encode: func(s string) (string, error) {
			return base32.StdEncoding.EncodeToString([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := base32.StdEncoding.DecodeString(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"base36": {
		Encode: func(s string) (string, error) {
			return base36.EncodeToStringUc([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := base36.DecodeString(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"base58": {
		Encode: func(s string) (string, error) {
			return base58.Encode([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := base58.Decode(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"base64": {
		Encode: func(s string) (string, error) {
			return base64.RawStdEncoding.EncodeToString([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := base64.RawStdEncoding.DecodeString(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"base64url": {
		Encode: func(s string) (string, error) {
			return base64.URLEncoding.EncodeToString([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := base64.URLEncoding.DecodeString(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"base85": {
		Encode: func(s string) (string, error) {
			return b85.Encode([]byte(s)), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := b85.Decode(s)
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		},
	},
	"ascii85": {
		Encode: func(s string) (string, error) {
			encoded := make([]byte, ascii85.MaxEncodedLen(len(s)))
			ascii85.Encode(encoded, []byte(s))
			return string(encoded), nil
		},
		Decode: func(s string) (string, error) {
			decoded := make([]byte, ascii85.MaxEncodedLen(len(s)))
			n, _, err := ascii85.Decode(decoded, []byte(s), true)
			if err != nil {
				return "", err
			}
			return string(decoded[:n]), nil
		},
	},
	"url": {
		Encode: func(s string) (string, error) {
			return url.QueryEscape(s), nil
		},
		Decode: func(s string) (string, error) {
			decoded, err := url.QueryUnescape(s)
			if err != nil {
				return "", err
			}
			return decoded, nil
		},
	},
	"morse": {
		Encode: func(s string) (string, error) {
			var buf strings.Builder
			for _, c := range s {
				for _, m := range morseEncoding {
					if !strings.EqualFold(m[0], string(c)) {
						continue
					}

					buf.WriteString(m[1] + " ")
					break
				}
			}
			return buf.String(), nil
		},
		Decode: func(s string) (string, error) {
			var buf strings.Builder
			parts := strings.Split(s, " ")
			for _, p := range parts {
				for _, m := range morseEncoding {
					if m[1] != p {
						continue
					}

					buf.WriteString(m[0])
					break
				}
			}
			return buf.String(), nil
		},
	},
}

var Hashes = map[string]func(string) string{
	"md5":     func(s string) string { return fmt.Sprintf("%x", md5.Sum([]byte(s))) },
	"sha1":    func(s string) string { return fmt.Sprintf("%x", sha1.Sum([]byte(s))) },
	"sha224":  func(s string) string { return fmt.Sprintf("%x", sha256.Sum224([]byte(s))) },
	"sha256":  func(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) },
	"sha384":  func(s string) string { return fmt.Sprintf("%x", sha512.Sum384([]byte(s))) },
	"sha512":  func(s string) string { return fmt.Sprintf("%x", sha512.Sum512([]byte(s))) },
	"blake2b": func(s string) string { return fmt.Sprintf("%x", blake2b.Sum256([]byte(s))) },
	"blake2s": func(s string) string { return fmt.Sprintf("%x", blake2s.Sum256([]byte(s))) },
	"adler32": func(s string) string { return fmt.Sprintf("%x", adler32.Checksum([]byte(s))) },
	"crc32":   func(s string) string { return fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(s))) },
	"crc64":   func(s string) string { return fmt.Sprintf("%x", crc64.Checksum([]byte(s), crc64.MakeTable(crc64.ISO))) },
	"fnv32":   func(s string) string { return fmt.Sprintf("%x", fnv.New32().Sum([]byte(s))) },
	"fnv64":   func(s string) string { return fmt.Sprintf("%x", fnv.New64().Sum([]byte(s))) },
	"fnv128":  func(s string) string { return fmt.Sprintf("%x", fnv.New128().Sum([]byte(s))) },
}

type Encryption struct {
	Encrypt func(string, string) (string, error)
	Decrypt func(string, string) (string, error)
}

var Encryptions = map[string]Encryption{
	"aes128": {
		Encrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Encrypt(s, 16) },
		Decrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Decrypt(s, 16) },
	},
	"aes192": {
		Encrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Encrypt(s, 24) },
		Decrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Decrypt(s, 24) },
	},
	"aes256": {
		Encrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Encrypt(s, 32) },
		Decrypt: func(s, key string) (string, error) { return encryption.AES{Key: key}.Decrypt(s, 32) },
	},
}

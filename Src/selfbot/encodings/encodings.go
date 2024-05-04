package encodings

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/ascii85"
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

	"github.com/darkwyrm/b85"
	"github.com/mr-tron/base58/base58"
	"github.com/multiformats/go-base32"
	"github.com/multiformats/go-base36"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

func Rot13(b rune) rune {
	if 'A' <= b && b <= 'Z' {
		return ((b-'A')+13)%26 + 'A'
	}
	if 'a' <= b && b <= 'z' {
		return ((b-'a')+13)%26 + 'a'
	}
	return b
}

func Hex(s string, decode bool) string {
	if decode {
		decoded, err := hex.DecodeString(s)
		if err != nil {
			return ""
		}
		return string(decoded)
	} else {
		return hex.EncodeToString([]byte(s))
	}
}

func Binary(s string, decode bool) string {
	if decode {
		var out []byte
		for i := 0; i+8 <= len(s); i += 8 {
			b, err := strconv.ParseUint(s[i:i+8], 2, 8)
			if err != nil {
				return "An error has occurred: " + err.Error()
			}
			out = append(out, byte(b))
		}
		return string(out)

	} else {
		var buf bytes.Buffer
		for i := 0; i < len(s); i++ {
			fmt.Fprintf(&buf, "%08b", s[i])
		}
		return buf.String()
	}
}

func Base32(s string, decode bool) string {
	b32 := base32.NewEncodingCI("abcdefghijklmnopqrstuvwxyz234567")
	if decode {
		decoded, err := b32.DecodeString(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return strings.ToUpper(b32.EncodeToString([]byte(s)))
	}
}

func Base36(s string, decode bool) string {
	if decode {
		decoded, err := base36.DecodeString(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return base36.EncodeToStringUc([]byte(s))
	}
}

func Base58(s string, decode bool) string {
	if decode {
		decoded, err := base58.Decode(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return base58.Encode([]byte(s))
	}
}

func Base64(s string, decode bool) string {
	if decode {
		decoded, err := base64.RawStdEncoding.DecodeString(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return base64.RawStdEncoding.EncodeToString([]byte(s))
	}
}

func Base64Url(s string, decode bool) string {
	if decode {
		decoded, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return base64.URLEncoding.EncodeToString([]byte(s))
	}
}

func Ascii85(s string, decode bool) string {
	if decode {
		decoded := make([]byte, ascii85.MaxEncodedLen(len(s)))
		n, _, err := ascii85.Decode(decoded, []byte(s), true)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded[:n])
	} else {
		encoded := make([]byte, ascii85.MaxEncodedLen(len(s)))
		ascii85.Encode(encoded, []byte(s))
		return string(encoded)
	}
}

func Base85(s string, decode bool) string {
	if decode {
		decoded, err := b85.Decode(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return string(decoded)
	} else {
		return b85.Encode([]byte(s))
	}
}

func URLEncode(s string, decode bool) string {
	if decode {
		decoded, err := url.QueryUnescape(s)
		if err != nil {
			return "An error has occurred: " + err.Error()
		}
		return decoded
	} else {
		return url.QueryEscape(s)
	}

}

func MorseCode(s string, decode bool) string {
	morseCode := [][]string{
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
	s = strings.ToUpper(s)
	if decode {
		for _, element := range morseCode {
			s = strings.ReplaceAll(s, element[1], element[0])
		}
	} else {
		for _, element := range morseCode {
			s = strings.ReplaceAll(s, element[0], element[1])
		}
	}
	return s
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func Sha224(s string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(s)))
}

func Sha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func Sha384(s string) string {
	return fmt.Sprintf("%x", sha512.Sum384([]byte(s)))
}

func Sha512(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func Blake2b(s string) string {
	return fmt.Sprintf("%x", blake2b.Sum256([]byte(s)))
}

func Blake2s(s string) string {
	return fmt.Sprintf("%x", blake2s.Sum256([]byte(s)))
}

func Adler32(s string) string {
	return fmt.Sprintf("%x", adler32.Checksum([]byte(s)))
}

func Crc32(s string) string {
	return fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(s)))
}

func Crc64(s string) string {
	return fmt.Sprintf("%x", crc64.Checksum([]byte(s), crc64.MakeTable(crc64.ISO)))
}

func Fnv32(s string) string {
	return fmt.Sprintf("%x", fnv.New32().Sum([]byte(s)))
}

func Fnv64(s string) string {
	return fmt.Sprintf("%x", fnv.New64().Sum([]byte(s)))
}

func Fnv128(s string) string {
	return fmt.Sprintf("%x", fnv.New128().Sum([]byte(s)))
}

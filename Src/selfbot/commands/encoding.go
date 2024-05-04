package commands

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/encodings"
	"github.com/LightningDev1/LB-Selfbot/selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type EncodingCommands struct{}

func (e EncodingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Encoding", "Encoding commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "rot13",
		Description: "Encodes or decodes text using ROT13",
		Usage:       "[p]rot13 <text>",
		Handler:     e.rot13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hex",
		Description: "Encodes or decodes text using hex",
		Usage:       "[p]hex <encode/decode> <text>",
		Handler:     e.hex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "binary",
		Description: "Encodes or decodes text using binary",
		Usage:       "[p]binary <encode/decode> <text>",
		Handler:     e.binary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base32",
		Description: "Encodes or decodes text using base32",
		Usage:       "[p]base32 <encode/decode> <text>",
		Handler:     e.base32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base36",
		Description: "Encodes or decodes text using base36",
		Usage:       "[p]base36 <encode/decode> <text>",
		Handler:     e.base36,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base58",
		Description: "Encodes or decodes text using base58",
		Usage:       "[p]base58 <encode/decode> <text>",
		Handler:     e.base58,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base64",
		Description: "Encodes or decodes text using base64",
		Usage:       "[p]base64 <encode/decode> <text>",
		Handler:     e.base64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base64url",
		Description: "Encodes or decodes text using URL safe base64",
		Usage:       "[p]base64url <encode/decode> <text>",
		Handler:     e.base64url,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ascii85",
		Description: "Encodes or decodes text using ascii85",
		Usage:       "[p]ascii85 <encode/decode> <text>",
		Handler:     e.ascii85,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base85",
		Description: "Encodes or decodes text using base85",
		Usage:       "[p]base85 <encode/decode> <text>",
		Handler:     e.base85,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "urlencode",
		Description: "Encodes or decodes text using URL encoding",
		Usage:       "[p]urlencode <encode/decode> <text>",
		Handler:     e.urlencode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "morsecode",
		Description: "Encodes or decodes text using morsecode",
		Usage:       "[p]morsecode <encode/decode> <text>",
		Handler:     e.morsecode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "md5",
		Description: "Hashes text using md5",
		Usage:       "[p]md5 <text>",
		Handler:     e.md5,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha1",
		Description: "Hashes text using sha1",
		Usage:       "[p]sha1 <text>",
		Handler:     e.sha1,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha224",
		Description: "Hashes text using sha224",
		Usage:       "[p]sha224 <text>",
		Handler:     e.sha224,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha256",
		Description: "Hashes text using sha256",
		Usage:       "[p]sha256 <text>",
		Handler:     e.sha256,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha384",
		Description: "Hashes text using sha384",
		Usage:       "[p]sha384 <text>",
		Handler:     e.sha384,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha512",
		Description: "Hashes text using sha512",
		Usage:       "[p]sha512 <text>",
		Handler:     e.sha512,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blake2b",
		Description: "Hashes text using blake2b",
		Usage:       "[p]blake2b <text>",
		Handler:     e.blake2b,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blake2s",
		Description: "Hashes text using blake2s",
		Usage:       "[p]blake2s <text>",
		Handler:     e.blake2s,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adler32",
		Description: "Hashes text using adler32",
		Usage:       "[p]adler32 <text>",
		Handler:     e.adler32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "crc32",
		Description: "Hashes text using crc32",
		Usage:       "[p]crc32 <text>",
		Handler:     e.crc32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "crc64",
		Description: "Hashes text using crc64",
		Usage:       "[p]crc64 <text>",
		Handler:     e.crc64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv32",
		Description: "Hashes text using fnv32",
		Usage:       "[p]fnv32 <text>",
		Handler:     e.fnv32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv64",
		Description: "Hashes text using fnv64",
		Usage:       "[p]fnv64 <text>",
		Handler:     e.fnv64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv128",
		Description: "Hashes text using fnv128",
		Usage:       "[p]fnv128 <text>",
		Handler:     e.fnv128,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aes128",
		Description: "Encrypt or decrypt text using AES-128",
		Usage:       "[p]aes128 <encrypt/decrypt> <\"key\"> <\"text\">",
		Handler:     e.aes128,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aes192",
		Description: "Encrypt or decrypt text using AES-192",
		Usage:       "[p]aes192 <encrypt/decrypt> <\"key\"> <\"text\">",
		Handler:     e.aes192,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aes256",
		Description: "Encrypt or decrypt text using AES-256",
		Usage:       "[p]aes256 <encrypt/decrypt> <\"key\"> <\"text\">",
		Handler:     e.aes256,
	})
}

func EncodingCommand(ctx *dgc.Ctx, encodeFunc func(string, bool) string) {
	encodeOrDecode := utils.Arguments["<encode/decode>"]
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	if encodeOrDecode == "encode" {
		_ = ctx.RespondText(encodeFunc(text, false))
	} else if encodeOrDecode == "decode" {
		_ = ctx.RespondText(encodeFunc(text, true))
	} else {
		_ = ctx.RespondText("Invalid option for <encode/decode>: " + encodeOrDecode)
	}
}

func (EncodingCommands) rot13(ctx *dgc.Ctx) {
	_ = ctx.RespondText(strings.Map(encodings.Rot13, ctx.Arguments.Raw()))
}

func (EncodingCommands) hex(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Hex)
}

func (EncodingCommands) binary(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Binary)
}

func (EncodingCommands) base32(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base32)
}

func (EncodingCommands) base36(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base36)
}

func (EncodingCommands) base58(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base58)
}

func (EncodingCommands) base64(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base64)
}

func (EncodingCommands) base64url(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base64Url)
}

func (EncodingCommands) ascii85(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Ascii85)
}

func (EncodingCommands) base85(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.Base85)
}

func (EncodingCommands) urlencode(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.URLEncode)
}

func (EncodingCommands) morsecode(ctx *dgc.Ctx) {
	EncodingCommand(ctx, encodings.MorseCode)
}

func (EncodingCommands) md5(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Md5(ctx.Arguments.Raw()))
}

func (EncodingCommands) sha1(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Sha1(ctx.Arguments.Raw()))
}

func (EncodingCommands) sha224(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Sha224(ctx.Arguments.Raw()))
}

func (EncodingCommands) sha256(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Sha256(ctx.Arguments.Raw()))
}

func (EncodingCommands) sha384(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Sha384(ctx.Arguments.Raw()))
}

func (EncodingCommands) sha512(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Sha512(ctx.Arguments.Raw()))
}

func (EncodingCommands) blake2b(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Blake2b(ctx.Arguments.Raw()))
}

func (EncodingCommands) blake2s(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Blake2s(ctx.Arguments.Raw()))
}

func (EncodingCommands) adler32(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Adler32(ctx.Arguments.Raw()))
}

func (EncodingCommands) crc32(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Crc32(ctx.Arguments.Raw()))
}

func (EncodingCommands) crc64(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Crc64(ctx.Arguments.Raw()))
}

func (EncodingCommands) fnv32(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Fnv32(ctx.Arguments.Raw()))
}

func (EncodingCommands) fnv64(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Fnv64(ctx.Arguments.Raw()))
}

func (EncodingCommands) fnv128(ctx *dgc.Ctx) {
	_ = ctx.RespondText(encodings.Fnv128(ctx.Arguments.Raw()))
}

func (EncodingCommands) aes128(ctx *dgc.Ctx) {
	encryptOrDecrypt := utils.Arguments["<encrypt/decrypt>"]
	key := utils.Arguments["<\"key\">"]
	text := utils.Arguments["<\"text\">"]
	if encryptOrDecrypt == "encrypt" {
		encrypted, err := encryption.AES{Key: key}.Encrypt(text, 16)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(encrypted)

	} else if encryptOrDecrypt == "decrypt" {
		decrypted, err := encryption.AES{Key: key}.Decrypt(text, 16)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(decrypted)
	} else {
		_ = ctx.RespondText("Invalid option for <encrypt/decrypt>: " + encryptOrDecrypt)
	}
}

func (EncodingCommands) aes192(ctx *dgc.Ctx) {
	encryptOrDecrypt := utils.Arguments["<encrypt/decrypt>"]
	key := utils.Arguments["<\"key\">"]
	text := utils.Arguments["<\"text\">"]
	if encryptOrDecrypt == "encrypt" {
		encrypted, err := encryption.AES{Key: key}.Encrypt(text, 24)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(encrypted)

	} else if encryptOrDecrypt == "decrypt" {
		decrypted, err := encryption.AES{Key: key}.Decrypt(text, 24)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(decrypted)
	} else {
		_ = ctx.RespondText("Invalid option for <encrypt/decrypt>: " + encryptOrDecrypt)
	}
}

func (EncodingCommands) aes256(ctx *dgc.Ctx) {
	encryptOrDecrypt := utils.Arguments["<encrypt/decrypt>"]
	key := utils.Arguments["<\"key\">"]
	text := utils.Arguments["<\"text\">"]
	if encryptOrDecrypt == "encrypt" {
		encrypted, err := encryption.AES{Key: key}.Encrypt(text, 32)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(encrypted)

	} else if encryptOrDecrypt == "decrypt" {
		decrypted, err := encryption.AES{Key: key}.Decrypt(text, 32)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		_ = ctx.RespondText(decrypted)
	} else {
		_ = ctx.RespondText("Invalid option for <encrypt/decrypt>: " + encryptOrDecrypt)
	}
}

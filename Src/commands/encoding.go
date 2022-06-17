package commands

import (
	"strings"

	"github.com/LightningDev/LB-Selfbot/encodings"
	"github.com/LightningDev/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type EncodingCommands struct{}

func (t EncodingCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "rot13",
		Description: "Encodes or decodes text using ROT13",
		Usage:       "[p]rot13 <text>",
		Handler:     t.rot13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hex",
		Description: "Encodes or decodes text using hex",
		Usage:       "[p]hex <encode/decode> <text>",
		Handler:     t.hex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "binary",
		Description: "Encodes or decodes text using binary",
		Usage:       "[p]binary <encode/decode> <text>",
		Handler:     t.binary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base32",
		Description: "Encodes or decodes text using base32",
		Usage:       "[p]base32 <encode/decode> <text>",
		Handler:     t.base32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base36",
		Description: "Encodes or decodes text using base36",
		Usage:       "[p]base36 <encode/decode> <text>",
		Handler:     t.base36,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base58",
		Description: "Encodes or decodes text using base58",
		Usage:       "[p]base58 <encode/decode> <text>",
		Handler:     t.base58,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base64",
		Description: "Encodes or decodes text using base64",
		Usage:       "[p]base64 <encode/decode> <text>",
		Handler:     t.base64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base64url",
		Description: "Encodes or decodes text using URL safe base64",
		Usage:       "[p]base64url <encode/decode> <text>",
		Handler:     t.base64url,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ascii85",
		Description: "Encodes or decodes text using ascii85",
		Usage:       "[p]ascii85 <encode/decode> <text>",
		Handler:     t.ascii85,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "base85",
		Description: "Encodes or decodes text using base85",
		Usage:       "[p]base85 <encode/decode> <text>",
		Handler:     t.base85,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "morsecode",
		Description: "Encodes or decodes text using morsecode",
		Usage:       "[p]morsecode <encode/decode> <text>",
		Handler:     t.morsecode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "md5",
		Description: "Hashes text using md5",
		Usage:       "[p]md5 <text>",
		Handler:     t.md5,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha1",
		Description: "Hashes text using sha1",
		Usage:       "[p]sha1 <text>",
		Handler:     t.sha1,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha224",
		Description: "Hashes text using sha224",
		Usage:       "[p]sha224 <text>",
		Handler:     t.sha224,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha256",
		Description: "Hashes text using sha256",
		Usage:       "[p]sha256 <text>",
		Handler:     t.sha256,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha384",
		Description: "Hashes text using sha384",
		Usage:       "[p]sha384 <text>",
		Handler:     t.sha384,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sha512",
		Description: "Hashes text using sha512",
		Usage:       "[p]sha512 <text>",
		Handler:     t.sha512,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blake2b",
		Description: "Hashes text using blake2b",
		Usage:       "[p]blake2b <text>",
		Handler:     t.blake2b,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blake2s",
		Description: "Hashes text using blake2s",
		Usage:       "[p]blake2s <text>",
		Handler:     t.blake2s,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adler32",
		Description: "Hashes text using adler32",
		Usage:       "[p]adler32 <text>",
		Handler:     t.adler32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "crc32",
		Description: "Hashes text using crc32",
		Usage:       "[p]crc32 <text>",
		Handler:     t.crc32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "crc64",
		Description: "Hashes text using crc64",
		Usage:       "[p]crc64 <text>",
		Handler:     t.crc64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv32",
		Description: "Hashes text using fnv32",
		Usage:       "[p]fnv32 <text>",
		Handler:     t.fnv32,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv64",
		Description: "Hashes text using fnv64",
		Usage:       "[p]fnv64 <text>",
		Handler:     t.fnv64,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fnv128",
		Description: "Hashes text using fnv128",
		Usage:       "[p]fnv128 <text>",
		Handler:     t.fnv128,
	})
}

func (t EncodingCommands) rot13(ctx *dgc.Ctx) {
	ctx.RespondText(strings.Map(encodings.Rot13, ctx.Arguments.Raw()))
}

func (t EncodingCommands) hex(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Hex)
}

func (t EncodingCommands) binary(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Binary)
}

func (t EncodingCommands) base32(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base32)
}

func (t EncodingCommands) base36(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base36)
}

func (t EncodingCommands) base58(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base58)
}

func (t EncodingCommands) base64(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base64)
}

func (t EncodingCommands) base64url(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base64Url)
}

func (t EncodingCommands) ascii85(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Ascii85)
}

func (t EncodingCommands) base85(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.Base85)
}

func (t EncodingCommands) morsecode(ctx *dgc.Ctx) {
	utils.EncodingCommand(ctx, encodings.MorseCode)
}

func (t EncodingCommands) md5(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Md5)
}

func (t EncodingCommands) sha1(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Sha1)
}

func (t EncodingCommands) sha224(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Sha224)
}

func (t EncodingCommands) sha256(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Sha256)
}

func (t EncodingCommands) sha384(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Sha384)
}

func (t EncodingCommands) sha512(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Sha512)
}

func (t EncodingCommands) blake2b(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Blake2b)
}

func (t EncodingCommands) blake2s(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Blake2s)
}

func (t EncodingCommands) adler32(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Adler32)
}

func (t EncodingCommands) crc32(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Crc32)
}

func (t EncodingCommands) crc64(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Crc64)
}

func (t EncodingCommands) fnv32(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Fnv32)
}

func (t EncodingCommands) fnv64(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Fnv64)
}

func (t EncodingCommands) fnv128(ctx *dgc.Ctx) {
	utils.HashCommand(ctx, encodings.Fnv128)
}

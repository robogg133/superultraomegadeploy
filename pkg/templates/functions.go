package templates

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/big"
	"strconv"
	"strings"
	"text/template"
)

const passswordCharSet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+=-"

var FuncMap template.FuncMap = template.FuncMap{
	"gen":          generate,
	"user_defined": func(kid string) any { return "foo bar" },
}

func generate(kind string) any {
	split := strings.Split(kind, "-")
	switch split[0] {
	case "random":
		if len(split) == 1 {
			return hex.EncodeToString(generateRandomBytes(32))
		}
		b := generateRandomBytes(mustAtoi(split[1]))

		if len(split) == 3 {
			return encode(split[2], b)
		}

		return hex.EncodeToString(b)
	case "password":
		if len(split) == 1 {
			return generatePassword(21)
		}
		n := mustAtoi(split[1])
		return generatePassword(n)
	case "ed25519":
		pk, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}

		if len(split) == 1 {
			return hex.EncodeToString(priv)
		}
		var b []byte

		switch split[1] {
		case "priv":
			b = priv
		case "pub":
			b = pk
		}

		if len(split) == 3 {
			return encode(split[2], b)
		}
		return hex.EncodeToString(b)

	default:
		return ""
	}
}

func generateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}
	return bytes
}

func generatePassword(length int) string {
	password := make([]byte, length)
	for i := range length {
		randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(passswordCharSet))))

		password[i] = passswordCharSet[randNum.Int64()]
	}

	return string(password)
}

func encode(kind string, data []byte) string {
	switch kind {
	case "rawbase64":
		return base64.RawStdEncoding.EncodeToString(data)
	case "base64":
		return base64.StdEncoding.EncodeToString(data)
	case "hex":
		return hex.EncodeToString(data)
	default:
		return hex.EncodeToString(data)
	}
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

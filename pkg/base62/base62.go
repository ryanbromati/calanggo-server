package base62

import (
	"errors"
	"strings"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = uint64(len(alphabet))

func Encode(id uint64) string {
	if id == 0 {
		return "0"
	}

	var sb strings.Builder
	for id > 0 {
		remainder := id % base
		sb.WriteByte(alphabet[remainder])
		id = id / base
	}

	return reverse(sb.String())
}

func Decode(token string) (uint64, error) {
	var id uint64 = 0

	for _, char := range token {
		index := strings.IndexRune(alphabet, char)
		if index == -1 {
			return 0, errors.New("caractere inválido no token base62")
		}
		id = id*base + uint64(index)
	}

	return id, nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

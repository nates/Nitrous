package main

import (
	"crypto/rand"
	"math/big"
)

func genCode() (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	code := make([]rune, 16)
	for i := range code {
		nBig, err := rand.Int(rand.Reader, big.NewInt(62))
		if err != nil {
			return "", err
		}
		code[i] = letters[nBig.Int64()]
	}
	return string(code), nil
}

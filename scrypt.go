package main

import (
	"crypto/rand"

	scrypt "golang.org/x/crypto/scrypt"
)

func scryptHashAndSalt(plaintext string) (hash []byte, salt []byte, err error) {
	salt, err = generateByteSliceToken(10)
	hash, err = scryptHash(plaintext, salt)
	return
}

func scryptHash(plaintext string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(plaintext), salt, 16384, 8, 1, 32)
}

func generateStringToken(length int) (token string, err error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const letterIdxBits = 6
	const letterIdxMask = 1<<letterIdxBits - 1
	tokenBytes := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes, err = generateByteSliceToken(bufferSize)
			if err != nil {
				return
			}
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			tokenBytes[i] = letterBytes[idx]
			i++
		}
	}
	token = string(tokenBytes)
	return
}

func generateByteSliceToken(length int) (token []byte, err error) {
	token = make([]byte, length)
	_, err = rand.Read(token)
	return
}

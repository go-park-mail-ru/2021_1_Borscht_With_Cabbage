package secure

import (
	"bytes"
	"context"
	"crypto/rand"

	"github.com/borscht/backend/utils/logger"
	"golang.org/x/crypto/argon2"
)

func GetSalt() []byte {
	salt := make([]byte, 8)
	rand.Read(salt)

	return salt
}

func HashPassword(ctx context.Context, salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)

	logger.UtilsLevel().DebugLog(ctx, logger.Fields{
		"password":      plainPassword,
		"salt":          salt,
		"hash password": hashedPass,
	})

	return append(salt, hashedPass...)
}

func CheckPassword(ctx context.Context, passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash)
	userPassHash := HashPassword(ctx, salt, plainPassword)

	logger.UtilsLevel().DebugLog(ctx, logger.Fields{
		"password": plainPassword,
		"expected": passHash,
		"received": userPassHash,
	})

	return bytes.Equal(userPassHash, passHash)
}

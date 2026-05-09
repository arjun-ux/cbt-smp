package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key untuk sign token, diambil dari Environment Variable (.env)
var jwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("KODE_RAHASIA_DEFAULT_YANG_SANGAT_KUAT_CBT_SMP_2026!")
	}
	return []byte(secret)
}()

// GenerateJWT membuat token baru untuk user yang berhasil login
func GenerateJWT(userID uint, username string, role string) (string, error) {
	// Menentukan expired time: Token berlaku 24 jam
	expirationTime := time.Now().Add(24 * time.Hour)

	// Membuat claims / payload
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      expirationTime.Unix(),
	}

	// Membuat token dengan metode signing HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT memvalidasi token dan mengembalikan claims
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Mengekstrak claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

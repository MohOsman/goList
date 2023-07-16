package authentication

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strings"
	"time"

	"goList/cache"
)

func generateSecretKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for a 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GenerateJWT(username string) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store the claims (payload) of the token
	claims := token.Claims.(jwt.MapClaims)

	// Set the claims
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Set token expiration time

	// Set the secret key used to sign the token
	secretKey, err := generateSecretKey()
	err = cache.Cache.Set("secert-key", secretKey)
	if err != nil {
		log.Printf("Error in cache")
		return "", err
	}

	if err != nil {
		log.Print("Error generating secret key")
		return "", err
	}

	// Generate the token string
	Skey, err := cache.Get("secert-key")
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(Skey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateAndExtractUsername(tokenString string) (string, error) {
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	tokenString = parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		skey, err := cache.Cache.Get("secert-key")
		if err != nil {
			log.Printf("Could not get secret key from cache")
			return nil, err
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return skey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	// Extract the username from the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	username, ok := claims["username"].(string)
	if !ok {
		log.Printf("Invalid token", err)
		return "", errors.New("invalid token")

	}
	return username, nil
}

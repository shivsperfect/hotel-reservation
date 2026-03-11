package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWT Authentication ---")
	token := c.Get("X-Api-Token")
	if token == "" {
		return fmt.Errorf("Unauthorized")
	}

	if err := parseToken(token); err != nil {
		return err
	}

	fmt.Println("token: ", token)
	return nil
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("Never print secret:", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("error parsing token: ", err)
		return fmt.Errorf("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("Unauthorized")
}

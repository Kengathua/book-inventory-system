package middlewares

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
)

// AuthRequiredMiddleware returns a Fiber middleware function that enforces authentication.
func AuthRequiredMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthRequired(c, db)
	}
}

// unauthorizedResponse sends an unauthorized response with the given message.
func unauthorizedResponse(c *fiber.Ctx, msg string) error {
	c.Status(fiber.StatusUnauthorized)
	return apis.HTTPStatusUnauthorizedResponse(c, errors.New(msg), msg, apis.ErrorResponseCode)
}

// parseAccessToken parses the access token.
func parseAccessToken(accessToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET_KEY")), nil
	})
}

// parseUserID parses the user ID from a string.
func parseUserID(userID string) (uuid.UUID, error) {
	return uuid.Parse(userID)
}

// AuthRequired is the authentication middleware that checks for a valid JWT token.
func AuthRequired(c *fiber.Ctx, db *gorm.DB) error {
	// Check if the user is already set in the context.
	if u, ok := c.Locals("user").(models.User); ok && u != (models.User{}) {
		// The user is already set, so proceed with the request.
		return nil
	}

	// Check the "Authorization" header for a JWT token.
	authorization := c.Get("Authorization")
	if authorization == "" {
		msg := "no token supplied"
		return unauthorizedResponse(c, msg)
	}

	tokenParts := strings.Split(authorization, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		msg := "invalid Authorization header"
		return unauthorizedResponse(c, msg)
	}

	accessToken := tokenParts[1]

	// Parse the access token.
	token, err := parseAccessToken(accessToken)
	if err != nil {
		msg := "invalid token"
		return unauthorizedResponse(c, msg)
	}

	// Look up the token in the database.
	var existingTokenInstance models.Token
	result := db.Where("access = ?", accessToken).First(&existingTokenInstance)
	if result.Error != nil || existingTokenInstance.Expired {
		msg := "token not found"
		return unauthorizedResponse(c, msg)
	}

	// Parse the issuer (user ID) from the token claims.
	claims := token.Claims.(*jwt.StandardClaims)
	id, err := parseUserID(claims.Issuer)
	if err != nil {
		msg := "error parsing user data"
		return unauthorizedResponse(c, msg)
	}

	// Fetch the user from the database and set it in the context.
	var user models.User
	db.Where("id = ?", id).First(&user)
	c.Locals("user", user)
	c.Next()

	return nil
}

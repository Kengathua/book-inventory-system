package apis

import (
	"fmt"
	"os"
	"time"

	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type RefreshRequestBody struct {
	Refresh string `json:"refresh"`
}

// GenerateAccessToken generates an access token for a user.
func GenerateAccessToken(user *models.User) (string, error) {
	// Create claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    *user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	// Generate encoded token
	access, err := claims.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return access, nil
}

// SaveNewToken saves a new token in the database and sets it as not expired.
func SaveNewToken(c *fiber.Ctx, db *gorm.DB, token *models.Token) error {
	token.Expired = false
	if result := db.Create(token); result.Error != nil {
		return result.Error
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token.Access,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return nil
}

func (h Handler) GetToken(c *fiber.Ctx) error {
	var user models.User
	var count int64
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing data", apis.ErrorResponseCode)
	}

	err = h.DB.Model(&models.User{}).Where("email = ?", data["email"]).Count(&count).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Invalid email or password", apis.ErrorResponseCode)
	}
	if count != 1 {
		return apis.HTTPStatusBadRequestResponse(c, fmt.Errorf("invalid email or password"), "Invalid email or password", apis.ErrorResponseCode)
	}
	result := h.DB.First(&user, "email = ?", data["email"])
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, result.Error, "Invalid email or password", apis.ErrorResponseCode)
	}

	isValid := user.CheckPasswordHarsh(data["password"])
	if !isValid {
		return apis.HTTPStatusBadRequestResponse(c, fmt.Errorf("could not match password to user"), "Invalid email or password", apis.ErrorResponseCode)
	}
	// Create claims and generate tokens
	access, err := GenerateAccessToken(&user)
	if err != nil {
		return apis.HTTPStatusInternalServerErrorResponse(c, err, "Invalid login", apis.ErrorResponseCode)
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{}).SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET_KEY")))
	if err != nil {
		return apis.HTTPStatusInternalServerErrorResponse(c, err, "Invalid login", apis.ErrorResponseCode)
	}

	token := &models.Token{Access: access, Refresh: refresh}
	if err := SaveNewToken(c, h.DB, token); err != nil {
		return apis.HTTPStatusInternalServerErrorResponse(c, err, "Invalid login", apis.ErrorResponseCode)
	}

	user.LastLogin = time.Now()
	if err := h.DB.Save(&user).Error; err != nil {
		return apis.HTTPStatusInternalServerErrorResponse(c, err, "Invalid login", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, token, "Successfully signed in", apis.SuccessResponseCode)
}

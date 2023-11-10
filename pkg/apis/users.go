package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type UserRequestBody struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

func (h Handler) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	if result := h.DB.Find(&users); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Users not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &users, "Successfully found users", apis.SuccessResponseCode)
}

func (h Handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if result := h.DB.First(&user, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "User not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &user, "Successfully found user", apis.SuccessResponseCode)
}

func (h Handler) AddUser(c *fiber.Ctx) error {
	body := UserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	var user models.User

	user.FullName = body.FullName
	user.Email = body.Email
	user.Password = body.Password
	user.UserType = body.UserType

	user.GeneratePasswordHarsh()

	err := h.DB.Create(&user).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error creating user", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusCreatedRequestResponse(c, &user, "Successfully created a new user", apis.SuccessResponseCode)
}

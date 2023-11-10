package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type AuthorRequestBody struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (h Handler) GetAuthors(c *fiber.Ctx) error {
	var authors []models.Author

	if result := h.DB.Preload("User").Find(&authors); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Authors not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &authors, "Successfully found authors", apis.SuccessResponseCode)
}

func (h Handler) AddAuthor(c *fiber.Ctx) error {
	body := AuthorRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	var author models.Author
	author.FullName = body.FullName
	author.Email = body.Email

	err := h.DB.Create(&author).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error creating author", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusCreatedRequestResponse(c, &author, "Successfully created a new author", apis.SuccessResponseCode)
}

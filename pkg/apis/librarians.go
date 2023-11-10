package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type LibrarianRequestBody struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (h Handler) GetLibrarians(c *fiber.Ctx) error {
	var librarians []models.Librarian

	if result := h.DB.Preload("User").Find(&librarians); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Librarians not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &librarians, "Successfully found librarians", apis.SuccessResponseCode)
}

func (h Handler) AddLibrarian(c *fiber.Ctx) error {
	body := LibrarianRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	var librarian models.Librarian
	librarian.FullName = body.FullName
	librarian.Email = body.Email

	err := h.DB.Create(&librarian).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error creating librarian", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusCreatedRequestResponse(c, &librarian, "Successfully created a new librarian", apis.SuccessResponseCode)
}

func (h Handler) LibrarianReviewBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var librarian models.Librarian

	if result := h.DB.First(&librarian, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Librarian not found", apis.ErrorResponseCode)
	}

	body := BookReviewRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing review book body", apis.ErrorResponseCode)
	}

	var book models.Book
	err := h.DB.First(&book, "id = ?", body.BookID).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error getting book", apis.ErrorResponseCode)
	}

	book.Status = "VERIFIED"
	book.IsLibrarianVerified = true
	err = h.DB.Save(&book).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error verifying book by librarian", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkRequestResponse(c, "Book successfuly verified by librarian", apis.SuccessResponseCode)
}

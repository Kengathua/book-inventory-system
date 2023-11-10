package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common"
	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type BookRequestBody struct {
	Title       string `json:"title"`
	AuthorID    string `json:"author_id"`
	ReleaseYear int    `json:"release_year"`
	Genre       string `json:"genre"`
	ISBNNumber  string `json:"isbn_number"`
	Quantity    int64  `json:"quantity"`
}

func (h Handler) GetBooks(c *fiber.Ctx) error {
	var books []models.Book

	if result := h.DB.Preload("Author.User").Find(&books); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Books not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &books, "Successfully found books", apis.SuccessResponseCode)
}

func (h Handler) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	if result := h.DB.Preload("Author.User").First(&book, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Book not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &book, "Successfully found book", apis.SuccessResponseCode)
}

func (h Handler) AddBook(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	userID, err := user.GetUserID_UUID()
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error getting signed in user id", apis.ErrorResponseCode)
	}

	// Parse request body
	body := BookRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	book := models.Book{
		BaseModel: common.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		Title:       body.Title,
		Genre:       body.Genre,
		AuthorID:    body.AuthorID,
		ReleaseYear: body.ReleaseYear,
		ISBNNumber:  body.ISBNNumber,
		Quantity:    body.Quantity,
	}

	err = h.DB.Create(&book).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error creating book", apis.ErrorResponseCode)
	}

	if result := h.DB.Preload("Author.User").First(&book, "id = ?", *book.ID); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Upated book not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusCreatedRequestResponse(c, &book, "Successfully created a new book", apis.SuccessResponseCode)
}

func (h Handler) UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	user := c.Locals("user").(models.User)
	userID, err := user.GetUserID_UUID()
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error getting signed in user id", apis.ErrorResponseCode)
	}

	if result := h.DB.First(&book, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Book not found", apis.ErrorResponseCode)
	}

	// Parse request body
	body := BookRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	book.UpdatedBy = userID
	book.Title = body.Title
	book.Genre = body.Genre
	book.AuthorID = body.AuthorID
	book.ReleaseYear = body.ReleaseYear
	book.ISBNNumber = body.ISBNNumber

	err = h.DB.Save(&book).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error updating book", apis.ErrorResponseCode)
	}

	if result := h.DB.Preload("Author.User").First(&book, "id = ?", *book.ID); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Upated book not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &book, "Successfully updated book", apis.SuccessResponseCode)
}

func (h Handler) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	if result := h.DB.First(&book, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Book not found", apis.ErrorResponseCode)
	}

	err := h.DB.Delete(&book).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error deleting book", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkRequestResponse(c, "Successfully deleted book", apis.SuccessResponseCode)
}

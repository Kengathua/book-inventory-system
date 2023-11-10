package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common/apis"
	"github.com/Kengathua/book-inventory-system/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type StoreKeeperRequestBody struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type BookReviewRequestBody struct {
	BookID string `json:"book_id"`
}

func (h Handler) GetStoreKeepers(c *fiber.Ctx) error {
	var storeKeepers []models.StoreKeeper

	if result := h.DB.Preload("User").Find(&storeKeepers); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "StoreKeepers not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &storeKeepers, "Successfully found storeKeepers", apis.SuccessResponseCode)
}

func (h Handler) GetStoreKeeper(c *fiber.Ctx) error {
	id := c.Params("id")
	var storeKeeper models.StoreKeeper

	if result := h.DB.Preload("Author.User").First(&storeKeeper, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Store keeper not found", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkDataRequestResponse(c, &storeKeeper, "Successfully found store keeper", apis.SuccessResponseCode)
}

func (h Handler) AddStoreKeeper(c *fiber.Ctx) error {
	body := StoreKeeperRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error parsing body", apis.ErrorResponseCode)
	}

	var storeKeeper models.StoreKeeper
	storeKeeper.FullName = body.FullName
	storeKeeper.Email = body.Email

	err := h.DB.Create(&storeKeeper).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error creating store keeper", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusCreatedRequestResponse(c, &storeKeeper, "Successfully created a new store keeper", apis.SuccessResponseCode)
}

func (h Handler) StoreKeeperReviewBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var storeKeeper models.StoreKeeper

	if result := h.DB.First(&storeKeeper, "id = ?", id); result.Error != nil {
		return apis.HTTPStatusNotFoundResponse(c, result.Error, "Store keeper not found", apis.ErrorResponseCode)
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

	book.Status = "APPROVED"
	book.IsKeeperVerified = true
	err = h.DB.Save(&book).Error
	if err != nil {
		return apis.HTTPStatusBadRequestResponse(c, err, "Error verifying book by store keeper", apis.ErrorResponseCode)
	}

	return apis.HTTPStatusOkRequestResponse(c, "Book successfuly verified by store keeper", apis.SuccessResponseCode)
}

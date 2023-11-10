package apis

import "github.com/gofiber/fiber/v2"

type BaseAPIResponse struct {
	Status          string `json:"status"`
	StatusCode      int    `json:"status_code" validate:"oneof=200 400 404"`
	ResponseMessage string `json:"responseMessage"`
	ResponseCode    int    `json:"responseCode"`
}

type APIOKResponse struct {
	BaseAPIResponse
}

type APIDataResponse struct {
	BaseAPIResponse
	Data interface{} `json:"data"`
}

type APIErrorResponse struct {
	BaseAPIResponse
	Error string `json:"error"`
}

func HTTPStatusOkRequestResponse(c *fiber.Ctx, msg string, responseCode int) error {
	response := APIOKResponse{}
	response.Status = "SUCCESS"
	response.StatusCode = fiber.StatusOK
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusOK).JSON(&response)
}

func HTTPStatusCreatedRequestResponse(c *fiber.Ctx, data interface{}, msg string, responseCode int) error {
	response := APIDataResponse{Data: data}
	response.Status = "SUCCESS"
	response.StatusCode = fiber.StatusCreated
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusCreated).JSON(&response)
}

func HTTPStatusOkDataRequestResponse(c *fiber.Ctx, data interface{}, msg string, responseCode int) error {
	response := APIDataResponse{Data: data}
	response.Status = "SUCCESS"
	response.StatusCode = fiber.StatusOK
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusOK).JSON(&response)
}

func HTTPStatusAcceptedResponse(c *fiber.Ctx, msg string, responseCode int) error {
	response := APIOKResponse{}
	response.Status = "SUCCESS"
	response.StatusCode = fiber.StatusAccepted
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusAccepted).JSON(&response)
}

func HTTPStatusBadRequestResponse(c *fiber.Ctx, err error, msg string, responseCode int) error {
	response := APIErrorResponse{Error: err.Error()}
	response.Status = "ERROR"
	response.StatusCode = fiber.StatusBadRequest
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusBadRequest).JSON(&response)
}

func HTTPStatusUnauthorizedResponse(c *fiber.Ctx, err error, msg string, responseCode int) error {
	response := APIErrorResponse{Error: err.Error()}
	response.Status = "ERROR"
	response.StatusCode = fiber.StatusUnauthorized
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusUnauthorized).JSON(&response)
}

func HTTPStatusNotFoundResponse(c *fiber.Ctx, err error, msg string, responseCode int) error {
	response := APIErrorResponse{Error: err.Error()}
	response.Status = "ERROR"
	response.StatusCode = fiber.StatusNotFound
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusNotFound).JSON(&response)
}

func HTTPStatusInternalServerErrorResponse(c *fiber.Ctx, err error, msg string, responseCode int) error {
	response := APIErrorResponse{Error: err.Error()}
	response.Status = "ERROR"
	response.StatusCode = fiber.StatusNotFound
	response.ResponseMessage = msg
	response.ResponseCode = responseCode

	return c.Status(fiber.StatusNotFound).JSON(&response)
}

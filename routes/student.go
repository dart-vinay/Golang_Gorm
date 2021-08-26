package routes

import (
	"Golang_Gorm/db"
	"Golang_Gorm/model"
	"Golang_Gorm/utils"
	"github.com/gofiber/fiber/v2"
)

func FetchSheetForStudentRequest(c *fiber.Ctx) error {
	headerInfo := model.FetchHeaderInfo(c)
	dbInstance := db.GORMDB()
	if response, err := headerInfo.FetchSheetsForStudent(dbInstance); err != nil {
		return handleErrorLogic(c, err)
	} else {
		return c.JSON(response)
	}
}

func PostUserComments(c *fiber.Ctx) error {
	headerInfo := model.FetchHeaderInfo(c)
	req := model.UserCommentPostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.ErrUnprocessableEntity)
	}
	req.HeaderInfo = headerInfo

	dbInstance := db.GORMDB()

	if err := req.PostUserComments(dbInstance); err != nil {
		return handleErrorLogic(c, err)
	} else {
		return c.JSON("Success!")
	}
	return nil
}

func handleErrorLogic(c *fiber.Ctx, err error) error {
	if err != nil {
		if err == utils.ErrBadRequest {
			return fiber.ErrBadRequest
		} else if err == utils.ErrInternal {
			return fiber.ErrInternalServerError
		} else if err == utils.ErrInvalidStudentId {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else if err == utils.ErrInvalidTeacherId {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else if err == utils.ErrUnknownUserType {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	return fiber.ErrInternalServerError
}

package routes

import (
	"Golang_Gorm/db"
	"Golang_Gorm/model"
	"github.com/gofiber/fiber/v2"
)

func FetchCommentsRelevantForTeacher(c *fiber.Ctx) error {
	headerInfo := model.FetchHeaderInfo(c)
	dbInstance := db.GORMDB()
	if response, err := headerInfo.FetchCommentsRelevantForTeacher(dbInstance); err != nil {
		return handleErrorLogic(c, err)
	} else {
		return c.JSON(response)
	}
}

func TakeActionForComment(c *fiber.Ctx) error {
	headerInfo := model.FetchHeaderInfo(c)
	req := model.ActionRequestForComments{}
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.ErrUnprocessableEntity)
	}
	req.HeaderInfo = headerInfo

	dbInstance := db.GORMDB()

	if err := req.TakeActionOnUserComments(dbInstance); err != nil {
		return handleErrorLogic(c, err)
	} else {
		return c.JSON("Success!")
	}
	return nil
}

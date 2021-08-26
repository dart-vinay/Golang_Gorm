package routes

import (
	"Golang_Gorm/db"
	"Golang_Gorm/model"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/common/log"
)

func FirstRoute(c *fiber.Ctx) error {
	log.Info("Entering")
	teacher := model.Teacher{
		ID: 1,
	}
	db := db.GORMDB()
	db.Where(&teacher).First(&teacher)
	return c.JSON(teacher)
}

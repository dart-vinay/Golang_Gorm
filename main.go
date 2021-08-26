package main

import (
	"Golang_Gorm/db"
	"Golang_Gorm/model"
	"Golang_Gorm/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/common/log"
)

// Overview:
// 1. Student can fetch all answer sheets with marks given for each answer for
// - Fetch sheet for student with subject name
// - Fetch all sheets for student : (Implemented)
// - Fetch all students in the class
// 2. Student can comment corresponding to a question for re-evaluation
// - POST : Comment for each question with question id and sheet id - (Implemented)
// - GET : All the comments posted by student ... return sheet and the question id with comment
// 3. Teacher can see all the comments logged for the sheets corrected by him/her and accept/reject
// - Get all sheets graded by teacher
// - Fetch all comments - (Implemented)
// - Resolve a comment by accepting or rejecting the comment: (Implemented)

func MigrateModels() {
	gormDB := db.GORMDB()

	err := gormDB.AutoMigrate(&model.Teacher{}, &model.Student{}, &model.Sheet{}, &model.Question{}, &model.SheetUserReference{}, &model.QuestionUserReference{}, &model.Comment{})
	if err != nil {
		log.Errorf("Error migrating models to Postgres %v", err)
		panic(err)
	}
}
func main() {

	app := fiber.New()

	db.InitDBConnection()
	MigrateModels()
	// uncomment this to populate data
	//data.PopulateData()
	app.Get("/v1/sheet/student", routes.FetchSheetForStudentRequest)       // FETCH ALL THE QUESTION SHEETS WITH MARKS AND COMMENTS AGAINST EACH QUESTION
	app.Post("/v1/sheet/student/comment", routes.PostUserComments)         // POST COMMENTS ON THE QUESTION
	app.Get("/v1/teacher/comment", routes.FetchCommentsRelevantForTeacher) // GET ALL RELEVANT COMMENTS FOR TEACHER
	app.Post("/v1/teacher/comment/action", routes.TakeActionForComment)    // ACTIONS TAKEN BY TEACHER FOR THE COMMENTS

	log.Fatal(app.Listen(":3000"))
}

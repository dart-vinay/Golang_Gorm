package data

import (
	"Golang_Gorm/db"
	"Golang_Gorm/model"
)

var students = []model.Student{
	model.Student{
		ID:   1,
		Name: "Vinay",
	},
	model.Student{
		ID:   2,
		Name: "Akshat",
	},
	model.Student{
		ID:   3,
		Name: "Saurabh",
	},
}

var teachers = []model.Teacher{
	model.Teacher{
		ID:                 1,
		Name:               "Teacher1",
		SubjectUndertaking: model.HINDI,
	},
	model.Teacher{
		ID:                 2,
		Name:               "Teacher2",
		SubjectUndertaking: model.ENGLISH,
	},
}

var sheets = []model.Sheet{
	model.Sheet{
		ID:          1,
		SubjectName: model.ENGLISH,
		//Questions: []model.Question{
		//	questionSet[0],
		//	questionSet[1],
		//},
	},
	model.Sheet{
		ID:          2,
		SubjectName: model.HINDI,
		//Questions: []model.Question{
		//	questionSet[2],
		//	questionSet[3],
		//},
	},
}

var questionSet = []model.Question{
	model.Question{
		ID:           1,
		SheetId:      1,
		QuestionText: "Question with ID 1",
		MaxMarks:     10,
	},
	model.Question{
		ID:           2,
		SheetId:      1,
		QuestionText: "Question with ID 2",
		MaxMarks:     20,
	},
	model.Question{
		ID:           3,
		SheetId:      2,
		QuestionText: "Question with ID 3",
		MaxMarks:     10,
	},
	model.Question{
		ID:           4,
		SheetId:      2,
		QuestionText: "Question with ID 4",
		MaxMarks:     20,
	},
}

var sheetsUserReference = []model.SheetUserReference{
	model.SheetUserReference{
		//ID:      1,
		SheetId: 1,
		UserId:  1,
	},
	model.SheetUserReference{
		//ID:      2,
		SheetId: 1,
		UserId:  2,
	},
	model.SheetUserReference{
		//ID:      3,
		SheetId: 2,
		UserId:  1,
	},
}

var questionUserReferences = []model.QuestionUserReference{
	model.QuestionUserReference{
		//ID: 1,
		UserId:       1,
		QuestionId:   1,
		MarksAlloted: 8,
		CorrectedBy:  2,
	},
	model.QuestionUserReference{
		//ID: 2,
		UserId:       1,
		QuestionId:   2,
		MarksAlloted: 15,
		CorrectedBy:  2,
	},
	model.QuestionUserReference{
		//ID: 3,
		UserId:       1,
		QuestionId:   3,
		MarksAlloted: 4,
		CorrectedBy:  1,
	},
	model.QuestionUserReference{
		//ID: 4,
		UserId:       1,
		QuestionId:   4,
		MarksAlloted: 10,
		CorrectedBy:  1,
	},
	model.QuestionUserReference{
		//ID: 5,
		UserId:       2,
		QuestionId:   1,
		MarksAlloted: 7,
		CorrectedBy:  2,
	},
	model.QuestionUserReference{
		//ID: 6,
		UserId:       2,
		QuestionId:   2,
		MarksAlloted: 10,
		CorrectedBy:  2,
	},
}

func PopulateData() {
	db := db.GORMDB()
	for _, student := range students {
		db.Save(&student)
	}

	for _, teacher := range teachers {
		db.Save(&teacher)
	}

	for _, sheet := range sheets {
		db.Save(&sheet)
	}
	for _, question := range questionSet {
		db.Save(&question)
	}
	for _, ref := range sheetsUserReference {
		db.Save(&ref)
	}

	for _, ref := range questionUserReferences {
		db.Save(&ref)
	}
}

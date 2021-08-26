package model

import (
	"gorm.io/gorm"
	"time"
)

// Generic Answer Sheet
type Sheet struct {
	gorm.Model
	ID          int `gorm:"primary_key"`
	SubjectName Subject
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	//Questions   []Question
}

//Generic Question in that sheet
type Question struct {
	gorm.Model
	ID           int `gorm:"primary_key"`
	QuestionText string
	MaxMarks     int
	SheetId      int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Sheet instance for user
type SheetUserReference struct {
	gorm.Model
	ID        int `gorm:"primary_key;AUTO_INCREMENT"`
	SheetId   int
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//
type QuestionUserReference struct {
	gorm.Model
	ID           int `gorm:"primary_key;AUTO_INCREMENT"`
	MarksAlloted int
	UserId       int
	CorrectedBy  int // teacher Id
	QuestionId   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Comment struct {
	gorm.Model
	ID                      int `gorm:"primary_key;AUTO_INCREMENT"`
	CommentText             string
	CommentStatus           CommentStatus
	UserId                  int
	QuestionUserReferenceId int
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt `gorm:"index"`
}

type CommentStatus string

const (
	REMARKED = CommentStatus("REMARKED")
	REJECT   = CommentStatus("REJECT")
	PENDING  = CommentStatus("PENDING")
)

func ValidStatusStrings() []string {
	return []string{string(PENDING), string(REJECT), string(REMARKED)}
}
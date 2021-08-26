package model

import (
	"Golang_Gorm/utils"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/common/log"
	"gorm.io/gorm"
	"time"
)

type Subject string

const (
	HINDI   = Subject("HINDI")
	ENGLISH = Subject("ENGLISH")
)

type Teacher struct {
	gorm.Model
	ID                 int `gorm:"primary_key"`
	Name               string
	SubjectUndertaking Subject
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type FetchTeacherRelevantCommentResponse struct {
	Comments []Comment `json:"comments"`
}

type ActionRequestForComments struct {
	HeaderInfo
	Actions []TeacherActionObject `json:"actions"`
}
type TeacherActionObject struct {
	QuestionRefId int    `json:"question_ref_id"`
	CommentId     int    `json:"comment_id"`
	CommentAction string `json:"comment_action"` // status for the comment
	ModifiedMarks int    `json:"modified_marks"`
}

func (req *HeaderInfo) FetchCommentsRelevantForTeacher(dbInstance *gorm.DB) (FetchTeacherRelevantCommentResponse, error) {
	comments := []Comment{}
	response := FetchTeacherRelevantCommentResponse{}
	teacherId := req.UserId
	validTeacherRequest, err := ValidateTeacher(dbInstance, teacherId)
	if err != nil {
		return response, err
	} else if validTeacherRequest.ID == 0 {
		return response, utils.ErrInvalidTeacherId
	}
	// Fetch All questions corrected by teacher
	questionsReferences := []QuestionUserReference{}
	dbInstance.Where("corrected_by=?", teacherId).Find(&questionsReferences)

	validQuestionRefId := []int{}
	for _, qRef := range questionsReferences {
		validQuestionRefId = append(validQuestionRefId, qRef.ID)
	}

	commentsMap, err := req.FetchCommentForQuestionRefs(validQuestionRefId, dbInstance)
	if err != nil {
		log.Errorf("Error occurred while fetching comments for question reference Ids %v req: %v", err, req)
		return response, err
	}
	for _, val := range commentsMap {
		comments = append(comments, val...)
	}
	return FetchTeacherRelevantCommentResponse{Comments: comments}, nil
}

func (req *ActionRequestForComments) TakeActionOnUserComments(dbInstance *gorm.DB) error {

	if req.UserId == 0 {
		return utils.ErrUnauthorized
	}
	questionRefId := []int{}
	actionsMap := make(map[int]TeacherActionObject)
	for _, action := range req.Actions {
		questionRefId = append(questionRefId, action.QuestionRefId)
		actionsMap[action.QuestionRefId] = action

	}
	questionRefId = utils.Unique(questionRefId)
	questionUserRefFromDB := []QuestionUserReference{}
	dbInstance.Find(&questionUserRefFromDB, "id IN ? AND corrected_by=?", questionRefId, req.UserId)

	for _, validQuestionRefs := range questionUserRefFromDB {
		refId := validQuestionRefs.ID
		commentStatus := strcase.ToScreamingSnake(actionsMap[refId].CommentAction)
		commentId := actionsMap[refId].CommentId
		if utils.ContainsString(ValidStatusStrings(), commentStatus) {

			if commentStatus == "REJECT" {
				dbInstance.Model(&Comment{}).Where("id=? AND question_user_reference_id=?", commentId, refId).Update("comment_status", commentStatus)
			} else if commentStatus == "REMARKED" {
				dbInstance.Model(&Comment{}).Where("id=? AND question_user_reference_id=?", commentId, refId).Update("comment_status", commentStatus)
				modifiedMarks := actionsMap[refId].ModifiedMarks
				dbInstance.Model(&QuestionUserReference{}).Where("id=? AND corrected_by=?", refId, req.UserId).Update("marks_alloted", modifiedMarks)
			} else {
				continue
			}
		} else {
			log.Errorf("Invalid comment action %v by teacher %v for question ref %v", commentStatus, req.UserId, refId)
		}
	}
	return nil
}

func ValidateTeacher(dbInstance *gorm.DB, teacherId int) (Teacher, error) {
	if teacherId == 0 {
		return Teacher{}, utils.ErrInvalidTeacherId
	}
	teacher := Teacher{}
	dbInstance.Where("id = ?", teacherId).Find(&teacher)
	return teacher, nil
}

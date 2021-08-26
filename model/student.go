package model

import (
	"Golang_Gorm/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/common/log"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type HeaderInfo struct {
	UserId   int      `json:"-"`
	UserType UserType `json:"-"`
}

type UserType string

const (
	STUDENT = UserType("STUDENT")
	TEACHER = UserType("TEACHER")
)

type Student struct {
	gorm.Model
	ID        int `gorm:"primary_key"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type FetchSheetForStudentResponse struct {
	SheetData []SheetData `json:"sheet_data"`
}

type UserCommentPostRequest struct {
	HeaderInfo
	UserComments []UserComments `json:"user_comments"`
}

type UserCommentPostResponse struct {
	UserComments []Comment `json:"user_comments"`
}

type UserComments struct {
	QuestionRefId int    `json:"question_ref_id"`
	Comment       string `json:"comment"`
}

type SheetData struct {
	SheetId   int
	Questions []QuestionDataResponseForUser
	Subject   Subject
}

type QuestionDataResponseForUser struct {
	QuestionId    int
	QuestionRefId int
	QuestionText  string
	CorrectedBy   int
	MaxMarks      int
	MarksAlloted  int
	Comments      []Comment
}

type UserAnsweredQuestions struct {
	UserId        int
	QuestionId    int
	QuestionRefId int
	TeacherId     int
	QuestionText  string
	SheetId       int
	MaxMarks      int
	MarksAlloted  int
}

func (req *HeaderInfo) FetchSheetsForStudent(dbInstance *gorm.DB) (FetchSheetForStudentResponse, error) {
	response := FetchSheetForStudentResponse{}

	// Validate Student
	studentId := req.UserId
	validStudentRequest, err := ValidateStudent(dbInstance, studentId)
	if err != nil {
		return response, err
	} else if validStudentRequest.ID == 0 {
		return response, utils.ErrInvalidStudentId
	}

	userAnsweredQuestions := []UserAnsweredQuestions{}
	rawSql := "Select questions.id as question_id, questions.sheet_id as sheet_id, questions.question_text as question_text, questions.max_marks as max_marks, question_user_references.user_id as user_id, question_user_references.corrected_by as teacher_id, question_user_references.id as question_ref_id, question_user_references.marks_alloted as marks_alloted from questions inner join question_user_references on questions.id = question_user_references.question_id where question_user_references.user_id=?"
	dbInstance.Raw(rawSql, studentId).Scan(&userAnsweredQuestions)
	questionRefIds := []int{}
	for _, userAnsweredQuestion := range userAnsweredQuestions {
		questionRefIds = append(questionRefIds, userAnsweredQuestion.QuestionRefId)
	}
	commentsMap, err := req.FetchCommentForQuestionRefs(questionRefIds, dbInstance)
	if err != nil {
		log.Errorf("Error fetching comments for question refs err : %v req %v", err, req)
		return response, err
	}
	sheetIds := []int{}
	for _, result := range userAnsweredQuestions {
		sheetIds = append(sheetIds, result.SheetId)
	}
	sheetIds = utils.Unique(sheetIds)
	sheetInfo := []Sheet{}
	dbInstance.Where(sheetIds).Find(&sheetInfo)

	sheetInfoMap := ConvertSheetArrayToMap(sheetInfo)

	sheetDataResponse := []SheetData{}
	for _, sheetId := range sheetIds {
		sheetData := SheetData{}
		questionsDataForResponse := []QuestionDataResponseForUser{}
		for _, questionsForUserRes := range userAnsweredQuestions {
			if questionsForUserRes.SheetId == sheetId {
				data := QuestionDataResponseForUser{
					QuestionId:    questionsForUserRes.QuestionId,
					QuestionText:  questionsForUserRes.QuestionText,
					QuestionRefId: questionsForUserRes.QuestionRefId,
					MarksAlloted:  questionsForUserRes.MarksAlloted,
					MaxMarks:      questionsForUserRes.MaxMarks,
					CorrectedBy:   questionsForUserRes.TeacherId,
					Comments:      commentsMap[questionsForUserRes.QuestionRefId],
				}
				questionsDataForResponse = append(questionsDataForResponse, data)
			}
		}
		if len(questionsDataForResponse) > 0 {
			sheetData.SheetId = sheetId
			sheetData.Subject = sheetInfoMap[sheetId].SubjectName
			sheetData.Questions = questionsDataForResponse

			sheetDataResponse = append(sheetDataResponse, sheetData)
		}
	}
	response.SheetData = sheetDataResponse
	return response, err
}

func (req *UserCommentPostRequest) PostUserComments(dbInstance *gorm.DB) error {
	//response := UserCommentPostResponse{}

	if req.UserId == 0 {
		return utils.ErrUnauthorized
	}
	questionRefId := []int{}
	//validQuestionRedId := []int{}
	commentsMap := make(map[int][]string)
	for _, commentObj := range req.UserComments {
		if commentObj.QuestionRefId != 0 {
			questionRefId = append(questionRefId, commentObj.QuestionRefId)
		}
		if _, ok := commentsMap[commentObj.QuestionRefId]; !ok {
			commentsMap[commentObj.QuestionRefId] = []string{strings.Trim(commentObj.Comment, " ")}
		} else {
			commentsMap[commentObj.QuestionRefId] = append(commentsMap[commentObj.QuestionRefId], strings.Trim(commentObj.Comment, " "))
		}
	}
	questionRefId = utils.Unique(questionRefId)
	questionUserRefFromDB := []QuestionUserReference{}
	dbInstance.Find(&questionUserRefFromDB, "id IN ? AND user_id=?", questionRefId, req.UserId)
	for _, validQuestionRef := range questionUserRefFromDB {
		refId := validQuestionRef.ID
		for _, commentText := range commentsMap[refId] {
			comment := Comment{
				UserId:                  req.UserId,
				CommentStatus:           PENDING,
				CommentText:             commentText,
				QuestionUserReferenceId: refId,
			}
			dbInstance.Save(&comment)
		}

	}
	return nil
}

func (req *HeaderInfo) FetchCommentForQuestionRefs(questionRefIds []int, dbInstance *gorm.DB) (map[int][]Comment, error) {
	comments := []Comment{}
	commentsMap := make(map[int][]Comment)
	dbInstance.Find(&comments, "question_user_reference_id IN ?", questionRefIds)
	for _, comment := range comments {
		if _, ok := commentsMap[comment.QuestionUserReferenceId]; !ok {
			commentsMap[comment.QuestionUserReferenceId] = []Comment{comment}
		} else {
			commentsMap[comment.QuestionUserReferenceId] = append(commentsMap[comment.QuestionUserReferenceId], comment)
		}
	}
	return commentsMap, nil
}

func ValidateStudent(dbInstance *gorm.DB, studentId int) (Student, error) {
	if studentId == 0 {
		return Student{}, utils.ErrInvalidStudentId
	}
	student := Student{}
	dbInstance.Where("id = ?", studentId).Find(&student)
	return student, nil
}

func FetchHeaderInfo(c *fiber.Ctx) HeaderInfo {
	header := HeaderInfo{}
	userId, _ := strconv.Atoi(string(c.Request().Header.Peek("X-USER-ID")))
	userType := strcase.ToScreamingSnake(string(c.Request().Header.Peek("X-USER-TYPE")))
	header.UserId = userId
	header.UserType = UserType(userType)

	return header
}

//func FetchSheetsForStudent(studentId int, dbInstance *gorm.DB) ([]FetchSheetForStudentResponse, error) {
//	response := []FetchSheetForStudentResponse{}
//	var err error
//
//	sheetReferencesResults := []SheetUserReference{}
//	questionReferencesResults := []QuestionUserReference{}
//
//	sheetUserRef := SheetUserReference{
//		UserId: studentId,
//	}
//	questionUserRef := QuestionUserReference{
//		UserId: studentId,
//	}
//	dbInstance.Where(&sheetUserRef).Find(&sheetReferencesResults)
//	dbInstance.Where(&questionUserRef).Find(&questionReferencesResults)
//
//	questionForUser := []int{}
//	//sheetsForUser := []int{}
//	marksAllotedMap := make(map[int]int)
//	for _, questionReferencesResult := range questionReferencesResults {
//		questionForUser = append(questionForUser, questionReferencesResult.QuestionId)
//		marksAllotedMap[questionReferencesResult.QuestionId] = questionReferencesResult.MarksAlloted
//	}
//
//	questions := []Question{}
//	dbInstance.Where(questionForUser).Find(&questions)
//	questionMap := ConvertQuestionArrayToMap(questions)
//
//	sheets := []Sheet{}
//	sheetsForUser := []int{}
//	for _, ques := range questions {
//		sheetsForUser = append(sheetsForUser, ques.SheetId)
//	}
//	sheetsForUser = utils.Unique(sheetsForUser)
//	dbInstance.Where(sheetsForUser).Find(&sheets)
//	//sheetMap := ConvertSheetArrayToMap(sheets)
//
//	sheetDataForResponse := []SheetData{}
//
//	for _, quesRef := range questionReferencesResults {
//		sheetId := questionMap[quesRef.QuestionId].SheetId
//		sheetRefId := int(0)
//		for _, sheetRef := range sheetReferencesResults {
//			if sheetRef.SheetId == sheetId {
//				sheetRefId = sheetRef.ID
//				break
//			}
//		}
//		sheetDataForResponse = append(sheetDataForResponse, SheetData{
//			SheetId:    sheetId,
//			SheetRefId: sheetRefId,
//		})
//	}
//	return response, err
//}

func ConvertQuestionArrayToMap(questions []Question) map[int]Question {
	res := make(map[int]Question)
	for _, ques := range questions {
		res[ques.ID] = ques
	}
	return res
}

func ConvertSheetArrayToMap(sheets []Sheet) map[int]Sheet {
	res := make(map[int]Sheet)
	for _, sheet := range sheets {
		res[sheet.ID] = sheet
	}
	return res
}

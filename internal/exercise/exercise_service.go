package exercise

import (
	"RestAPI/internal/domain"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseService struct {
	db *gorm.DB
}

func NewExerciseService(database *gorm.DB) *ExerciseService {
	return &ExerciseService{
		db: database,
	}
}

func (ex ExerciseService) GetExercise(ctx *gin.Context) {
	paramID := ctx.Param("id")
	//fmt.Println(paramID)
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	//fmt.Println(id)

	var exercise domain.Exercise
	err = ex.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, exercise)
}

func (ex ExerciseService) GetUserScore(ctx *gin.Context) {
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = ex.db.Where("exercise_id = ? AND user_id = ?", exerciseID, userID).Find(&answers).Error
	if err != nil {
		ctx.JSON(200, gin.H{
			"score": 0,
		})
		return
	}
	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score int
	for _, question := range exercise.Questions {
		if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
			score += question.Score
		}
	}
	ctx.JSON(200, gin.H{
		"score": score,
	})
}

//tugas

func (ex ExerciseService) CreateExercise(ctx *gin.Context) {
	var excReq domain.ExerciseRequest
	err := ctx.ShouldBind(&excReq)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	exc := domain.Exercise{
		Title:       excReq.Title,
		Description: excReq.Description,
	}

	if err := ex.db.Create(&exc).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when creating exercise",
		})
		return
	}
	ctx.JSON(201, exc)
}

func (ex ExerciseService) CreateQuestion(ctx *gin.Context) {
	var queReq domain.QuestionRequest
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	err = ctx.ShouldBind(&queReq)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	var exc domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Take(&exc).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "exercise id is not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))

	que := domain.Question{
		ExerciseID:    exerciseID,
		Body:          queReq.Body,
		OptionA:       queReq.OptionA,
		OptionB:       queReq.OptionB,
		OptionC:       queReq.OptionC,
		OptionD:       queReq.OptionD,
		CorrectAnswer: queReq.CorrectAnswer,
		Score:         100,
		CreatorID:     userID,
	}

	if err := ex.db.Create(&que).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when creating user",
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "successful operation",
	})

}

func (ex ExerciseService) CreateAnswer(ctx *gin.Context) {
	var ansReq domain.AnswerRequest
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	paramQuestionID := ctx.Param("qid")
	questionID, err := strconv.Atoi(paramQuestionID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid question id",
		})
		return
	}

	err = ctx.ShouldBind(&ansReq)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	var que domain.Question
	err = ex.db.Where("id = ? AND exercise_id = ?", questionID, exerciseID).Take(&que).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "question not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))

	ans := domain.Answer{
		ExerciseID: exerciseID,
		QuestionID: questionID,
		UserID:     userID,
		Answer:     ansReq.Answer,
	}

	if err := ex.db.Create(&ans).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when submitting answer",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "successful operation",
	})

}

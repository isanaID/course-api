package usecase

import (
	"course-api/internal/domain"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseUsecase struct {
	db *gorm.DB
}

func NewExerciseUsecase(db *gorm.DB) *ExerciseUsecase {
	return &ExerciseUsecase{
		db: db,
	}
}

func (exerUsecase *ExerciseUsecase) CreateExercise(c *gin.Context) {
	var exercise domain.Exercise
	c.BindJSON(&exercise)
	exerUsecase.db.Create(&exercise)
	c.JSON(200, gin.H{
		"message": "Exercise created successfully!",
	})
}

func (exerUsecase *ExerciseUsecase) GetExercises(c *gin.Context) {
	var exercises []domain.Exercise
	exerUsecase.db.Find(&exercises)
	c.JSON(200, gin.H{
		"exercises": exercises,
	})
}

func (exerUsecase *ExerciseUsecase) GetExercise(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}

	var exercise domain.Exercise
	err = exerUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	c.JSON(200, exercise)
}

func (exerUsecase *ExerciseUsecase) UpdateExercise(c *gin.Context) {
	var exercise domain.Exercise
	id := c.Param("id")
	exerUsecase.db.First(&exercise, id)
	c.BindJSON(&exercise)
	exerUsecase.db.Save(&exercise)
	c.JSON(200, gin.H{
		"message": "Exercise updated successfully!",
	})
}

func (exerUsecase *ExerciseUsecase) DeleteExercise(c *gin.Context) {
	var exercise domain.Exercise
	id := c.Param("id")
	exerUsecase.db.First(&exercise, id)
	exerUsecase.db.Delete(&exercise)
	c.JSON(200, gin.H{
		"message": "Exercise deleted successfully!",
	})
}

func (exerUsecase *ExerciseUsecase) CreateQuestion(c *gin.Context) {
	type Question struct {
		Body          string `json:"body"`
		OptionA       string `json:"option_a"`
		OptionB       string `json:"option_b"`
		OptionC       string `json:"option_c"`
		OptionD       string `json:"option_d"`
		CorrectAnswer string `json:"correct_answer"`
	}

	var question Question
	c.BindJSON(&question)
	var exercise domain.Exercise
	exerUsecase.db.First(&exercise, c.Param("id"))
	exercise.Questions = append(exercise.Questions, domain.Question{
		Body:          question.Body,
		OptionA:       question.OptionA,
		OptionB:       question.OptionB,
		OptionC:       question.OptionC,
		OptionD:       question.OptionD,
		CorrectAnswer: question.CorrectAnswer,
	})
	exerUsecase.db.Save(&exercise)
	c.JSON(201, gin.H{
		"message": "Question created successfully!",
	})

	if question.Body == "" || question.OptionA == "" || question.OptionB == "" || question.OptionC == "" || question.OptionD == "" || question.CorrectAnswer == "" {
		c.JSON(400, gin.H{
			"message": "Question body, options and correct answer are required",
		})
		return
	}
}

func (exerUsecase *ExerciseUsecase) CreateAnswer(c *gin.Context) {
	var answer domain.Answer
	c.BindJSON(&answer)
	var exercise domain.Exercise
	exerUsecase.db.First(&exercise, c.Param("id"))
	var question domain.Question
	exerUsecase.db.First(&question, answer.QuestionID)
	var scoreCount ScoreCount
	exerUsecase.db.First(&scoreCount, "exercise_id = ?", c.Param("id"))
	if answer.Answer == question.CorrectAnswer {
		scoreCount.Inc(1)
	}
	exerUsecase.db.Save(&scoreCount)
	c.JSON(201, gin.H{
		"message": "Answer created successfully!",
	})
}

func (exerUsecase *ExerciseUsecase) CalculateScore(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}

	var exercise domain.Exercise
	err = exerUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = exerUsecase.db.Where("exercise_id = ? AND user_id = ?", id, userID).Find(&answers).Error
	if err != nil || len(answers) == 0 {
		c.JSON(200, map[string]interface{}{
			"score": 0,
		})
		return
	}
	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score ScoreCount
	wg := new(sync.WaitGroup)
	for _, question := range exercise.Questions {
		newQuestion := question
		wg.Add(1)
		go func() {
			defer wg.Done()
			if strings.EqualFold(newQuestion.CorrectAnswer, mapQA[newQuestion.ID].Answer) {
				score.Inc(newQuestion.Score)
			}
		}()
	}
	wg.Wait()
	c.JSON(200, map[string]interface{}{
		"score": score.score,
	})
}

type ScoreCount struct {
	score int
	mu    sync.Mutex
}

func (sc *ScoreCount) Inc(value int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.score += value
}

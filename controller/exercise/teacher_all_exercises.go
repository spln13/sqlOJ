package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type TeacherAllExercise struct {
	List []TeacherOneExercise `json:"list"`
	utils.Response
}

type TeacherOneExercise struct {
	Answer        string `json:"answer"`
	Description   string `json:"description"`
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherID   int64  `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
}

func TeacherGetAllExercises(context *gin.Context) {
	allExercisesList, err := model.NewExerciseContentFlow().GetAllExercise()
	if err != nil {
		context.JSON(http.StatusOK, TeacherAllExercise{
			List:     nil,
			Response: utils.NewCommonResponse(0, ""),
		})
		return
	}
	var teacherAllExerciseList []TeacherOneExercise
	for _, exercise := range allExercisesList {
		oneExercise := TeacherOneExercise{
			Answer:        exercise.Answer,
			Description:   exercise.Description,
			ExerciseID:    exercise.ID,
			ExerciseName:  exercise.Name,
			Grade:         exercise.Grade,
			PassCount:     exercise.PassCount,
			PublisherID:   exercise.PublisherID,
			PublisherName: exercise.PublisherName,
			PublisherType: exercise.PublisherType,
			SubmitCount:   exercise.SubmitCount,
		}
		teacherAllExerciseList = append(teacherAllExerciseList, oneExercise)
	}
	context.JSON(http.StatusOK, TeacherAllExercise{
		List:     teacherAllExerciseList,
		Response: utils.NewCommonResponse(0, ""),
	})
}

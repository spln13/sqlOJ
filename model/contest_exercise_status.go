package model

import (
	"sync"
	"time"
)

type ContestExerciseStatus struct {
	ID         int64
	ContestID  int64
	ExerciseID int64
	UserID     int64
	UserType   int64
	Status     int64 // 1->ac; 2->wa; 3->re
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ContestExerciseStatusFlow struct {
}

var (
	contestExerciseStatusFlowDAO  *ContestExerciseStatusFlow
	contestExerciseStatusFlowOnce sync.Once
)

func NewContestExerciseStatusFlow() *ContestExerciseStatusFlow {
	contestExerciseStatusFlowOnce.Do(func() {
		contestExerciseStatusFlowDAO = new(ContestExerciseStatusFlow)
	})
	return contestExerciseStatusFlowDAO
}

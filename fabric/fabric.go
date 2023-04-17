package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

import (
	"encoding/json"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type Submission struct {
	UserID     string `json:"user_id"`
	UserType   string `json:"user_type"`
	ExerciseID string `json:"exercise_id"`
	ContestID  string `json:"contest_id"`
	Status     string `json:"status"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Submission
}

// CreateSubmission adds a new submission to the world state with given details
func (s *SmartContract) CreateSubmission(ctx contractapi.TransactionContextInterface, submissionNumber string, userID string, userType string, exerciseID string, contestID string, status string) error {
	submission := Submission{
		UserID:     userID,
		UserType:   userType,
		ExerciseID: exerciseID,
		ContestID:  contestID,
		Status:     status,
	}

	submissionAsBytes, _ := json.Marshal(submission)

	return ctx.GetStub().PutState(submissionNumber, submissionAsBytes)
}

// QuerySubmission returns the submission stored in the world state with given id
func (s *SmartContract) QuerySubmission(ctx contractapi.TransactionContextInterface, submissionNumber string) (*Submission, error) {
	submissionAsBytes, err := ctx.GetStub().GetState(submissionNumber)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if submissionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", submissionNumber)
	}

	submission := new(Submission)
	_ = json.Unmarshal(submissionAsBytes, submission)

	return submission, nil
}

// ChangeSubmissionStatus updates the status field of submission with given id in world state
func (s *SmartContract) ChangeSubmissionStatus(ctx contractapi.TransactionContextInterface, submissionNumber string, newStatus string) error {
	submission, err := s.QuerySubmission(ctx, submissionNumber)

	if err != nil {
		return err
	}

	submission.Status = newStatus

	submissionAsBytes, _ := json.Marshal(submission)

	return ctx.GetStub().PutState(submissionNumber, submissionAsBytes)
}

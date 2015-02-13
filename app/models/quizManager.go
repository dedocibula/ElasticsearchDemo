package models

import (
	"time"
)

const (
	index       = "dba"
	quizType    = "question"
	rankingType = "ranking"
)

type QuizManager struct {
	em *ELKManager
}

func NewQuizManager() *QuizManager {
	rm := NewResourceManager()
	em := NewELKManager(rm)
	return &QuizManager{em: em}
}

func (q *QuizManager) Dispose() {
	q.em.Dispose()
	q.em = nil
}

func (q QuizManager) Validate(attempt Attempt) Result {
	result, err := q.em.LiteralSearchELK(index, quizType)
	if err != nil {
		return q.buildResult(false, err.Error())
	} else if result != attempt.Answer {
		return q.buildResult(false,
			"Sorry, your answer wasn't quite correct. Please try again.")
	} else {
		record := q.createELKRecord(attempt.Name)
		return q.submitELKRecord(record)
	}
}

func (q QuizManager) buildResult(correct bool, message string) Result {
	return Result{
		Ok:      correct,
		Message: message,
	}
}

func (q QuizManager) createELKRecord(nickname string) ELKRecord {
	return ELKRecord{
		Nickname:  nickname,
		Timestamp: time.Now(),
	}
}

func (q QuizManager) submitELKRecord(record ELKRecord) Result {
	err := q.em.RecordSuccessELK(index, rankingType, record)
	if err != nil {
		return q.buildResult(false, err.Error())
	} else {
		return q.buildResult(true, "That's the correct answer. Well done.")
	}
}
